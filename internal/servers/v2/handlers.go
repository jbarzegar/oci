package serverv2

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jbarzegar/oci/internal/blobstorage"
)

// handler stores http handlers
// passes dependencies to handlers from the root subapp
type handler struct {
	storageClient blobstorage.Storer
}

// func handlePullManifest(c fiber.Ctx) error {
// 	invalid := []Error{}

// 	// validate each item then append any potential invalid
// 	// params to a slice. Allows for fixing all potential issues on first pass
// 	// rather than failing-fast each time something happens
// 	name := c.Params("name", "")
// 	nameValid := ociv2.ValidateManifestName(name)
// 	if !nameValid {
// 		invalid = append(invalid,
// 			serverError(ERR_NAME_INVALID, "param name invalid", ""))
// 	}

// 	ref := c.Params("reference", "")
// 	refValid := ociv2.ValidateReference(ref)
// 	if !refValid {
// 		invalid = append(invalid,
// 			serverError(ERR_MANIFEST_INVALID, "invalid manifest", ""))
// 	}

// 	if len(invalid) > 0 {
// 		return handleErrorResponse(c, 401, invalid...)
// 	}

// 	return nil
// }

func (h *handler) BlobUploads(c fiber.Ctx) error {
	log.WithContext(c).Info("start blob upload")

	errs := []Error{}

	n := c.Params("name", "")
	if n == "" {
		errs = append(errs,
			serverError(ERR_BLOB_UPLOAD_INVALID, "name param not set", ""))

	}

	log.Debugw("thing", "name", n)

	if len(errs) > 1 {
		return handleErrorResponse(c, 401, errs...)
	}

	return handleErrorResponse(c, 401,
		serverError(ERR_NOT_IMPLEMENTED,
			"uploads not implemented yet", map[string]string{"name": n},
		))
}
