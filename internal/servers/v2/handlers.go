package serverv2

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/jbarzegar/oci/internal/blobstorage"
	"github.com/jbarzegar/oci/internal/manifest"
)

// handler stores http handlers
// passes dependencies to handlers from the root subapp
type handler struct {
	storageClient blobstorage.Storer
}

// BlobExists checks if a given name & digest exists in the
// storage client and returns a subsiquent http status if found:
// 200 not found: 404
func (h *handler) BlobExists(c fiber.Ctx) error {
	digest := c.Params("digest")
	if digest == "" {
		return c.SendStatus(404)
	}
	name := c.Params("name")
	if name == "" {
		return c.SendStatus(404)
	}

	size, exists, err := h.storageClient.BlobInfo(c.Context(), name, digest)
	if err != nil {
		log.WithContext(c).Errorw("couldn't get blob info", "error", err)
		return c.SendStatus(404)
	}

	if !exists {
		return c.SendStatus(404)
	}
	c.Response().Header.Add("Accept-Ranges", "bytes")
	// response.Header().Set("Content-Type", resolveBlobResponseMediaType(imgStore, name, digest, rh.c.Log))
	// TODO: understand this ^
	c.Response().Header.Add("Content-Type", "application/octet-stream")
	c.Response().Header.Add("Content-Length", strconv.FormatInt(size, 10))
	// TODO: This ->
	// response.Header().Set(constants.DistContentDigestKey, digest.String())

	return c.SendStatus(200)
}

// BlobRefClose handles when a blob upload is closed by the client
// this is where a blob is actually written to a storage solution
func (h *handler) BlobRefClose(c fiber.Ctx) error {
	queries := c.Queries()

	ref := c.Params("reference")
	if ref == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	name := c.Params("name")
	if name == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if digest, ok := queries["digest"]; ok {
		log.WithContext(c).Infow("closing upload stream",
			"name", name,
			"reference", ref,
			"digest", digest,
		)

		writer, err := h.storageClient.GetWriterByName(c.Context(), name)
		defer writer.Flush()
		if err != nil {
			log.Error("couldn't find writer by name", "name", name)
			return c.Status(404).Send([]byte{})
		}
		err = writer.Write(c.Context(), digest)
		if err != nil {
			log.Errorw("failed to write image", "error", err, "name", name)
			return c.Status(400).Send([]byte(err.Error()))
		}

	} else {
		return c.SendStatus(404)
	}

	return c.SendStatus(201)
}

// BlobUpload handles uploading a given blob/layer of an image
// given a reference and name, BlobUpload will write that part to
// a stream, return the content length and prepare the blob to be
// written. The method doesn't actually write it to disk however
func (h *handler) BlobUpload(c fiber.Ctx) error {
	ref := c.Params("reference", "")
	uid, err := uuid.Parse(ref)
	if err != nil {
		return handleErrorResponse(c, 400,
			serverError(ERR_BLOB_UPLOAD_INVALID,
				"could not parse Session ID", map[string]string{"reference": ref},
			))
	}
	name := c.Params("name")

	contentRange := 0
	reqHeaders := c.Req().GetHeaders()
	cr, ok := reqHeaders["Content-Range"]
	if !ok {
		log.WithContext(c).Warnw("No Content-Range header found, wtf do I need to do") // "headers", c.GetHeaders(),
		// "reqHeaders", reqHeaders,

	} else {
		log.WithContext(c).Warnw("GOT CONTENT RANGE", "CONTENT-RANGE", cr)
	}

	if !c.Req().Is("application/octet-stream") {
		return handleErrorResponse(c, 400,
			serverError(ERR_BLOB_UPLOAD_INVALID,
				"Request must be an octet stream", reqHeaders["Content-Type"],
			))
	}
	// Fetch the writer for a given uid
	w, err := h.storageClient.GetWriterByName(c.Context(), name)
	if err != nil {
		log.WithContext(c).Fatalw("failed to get writer",
			"error", err,
		)
		return handleErrorResponse(c, 400,
			serverError(ERR_BLOB_UPLOAD_UNKNOWN, "failed to get writer", err),
		)
	}
	// Append the bytes from the body
	clen, err := w.AppendPart(c.Context(), uid, c.Body())
	if err != nil {
		log.WithContext(c).Fatalw("failed to write part",
			"error", err,
		)
		return handleErrorResponse(c, 400,
			serverError(ERR_BLOB_UPLOAD_UNKNOWN, "failed to write part", err.Error()),
		)
	}
	// generate current range
	c.Response().Header.Add("Range", fmt.Sprintf("%v-%v", contentRange, clen))
	return c.Status(202).Send([]byte{})
}

// BlobUploadLocation is the entrypoint for a chunked blob upload.
// It will respond with a accepted response range, and UUID of
// where the blob must be uploaded
func (h *handler) BlobUploadLocation(c fiber.Ctx) error {
	n := c.Params("name", "")
	if n == "" {
		return handleErrorResponse(c, 404,
			serverError(ERR_BLOB_UPLOAD_INVALID, "name param not set", ""),
		)
	}

	// Currently do not support full upload
	if c.Query("digest", "") != "" {
		return c.SendStatus(fiber.ErrMethodNotAllowed.Code)
	}

	contentLength := c.Request().Header.ContentLength()
	if contentLength == 0 {
		log.Infow("blob upload will be chunked")
	} else if contentLength < 0 {
		return handleErrorResponse(
			c, 401,
			serverError(ERR_BLOB_UPLOAD_INVALID,
				"content-length cannot be negagtive",
				map[string]int{"Content-Length": contentLength},
			))
	} else {
		log.Infow("content length is beeg is being pushed already?",
			"Content-Length", contentLength,
			"Body", len(c.Body()),
		)

	}

	// generate session ID (uuid)
	uid := uuid.New()
	// We'll create the writer however we have nothing to write
	// as of yet this will be done when we PATCH
	_, err := h.storageClient.CreateWriter(c.Context(), n, uid)
	if err != nil {
		return handleErrorResponse(
			c, 400,
			serverError(ERR_BLOB_UPLOAD_UNKNOWN,
				"failed to create writer",
				map[string]any{
					"err": err,
				},
			))
	}

	c.Response().Header.Add("Location", uid.String())
	c.Response().Header.Add("Range", "0-0")
	return c.SendStatus(fiber.StatusAccepted)
}

// UploadManifest
func (h *handler) UploadManifest(c fiber.Ctx) error {
	name := c.Params("name")
	reference := c.Params("reference")
	fmt.Println("tags", c.Req().OriginalURL())
	parsedManifest, err := manifest.UnmarshalV2(c.Body())
	if err != nil {
		if errors.Is(err, manifest.ErrManifestInvalid) {
			return handleErrorResponse(c,
				400,
				serverError(ERR_MANIFEST_INVALID, "manifest unparsable", err),
			)
		}
		return err
	}
	err = h.storageClient.WriteManifest(
		c.Context(),
		name,
		reference,
		parsedManifest,
	)
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("%v/v2/%v/manifests/%v", c.BaseURL(), name, reference)
	c.Response().Header.Add("Location", loc)
	c.Response().Header.Add("Docker-Content-Digest", parsedManifest.Config.Digest)
	c.Response().Header.Add("OCI-Tag", reference)
	return c.SendStatus(201)
}
