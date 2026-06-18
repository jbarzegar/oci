package serverv2

import (
	"errors"
	"fmt"

	"github.com/codingconcepts/env"
	"github.com/gofiber/fiber/v3"
	"github.com/jbarzegar/oci/internal/blobstorage"
	"github.com/jbarzegar/oci/internal/manifest"
)

type V2ServerConfig struct {
	BlobStorageBucketName   string `env:"BLOB_STORAGE_BUCKET_NAME" required:"true"`
	BlobStorageCreateBucket bool   `env:"BLOB_STORAGE_CREATE_BUCKET" default:"false"`
}

func loadServerConfig() (*V2ServerConfig, error) {
	serverConfig := V2ServerConfig{}
	if err := env.Set(&serverConfig); err != nil {
		return nil, err
	}

	return &serverConfig, nil
}

func New() (*fiber.App, error) {
	cfg, err := loadServerConfig()
	if err != nil {
		return nil, err
	}

	app := fiber.New()

	storageClient, err := blobstorage.NewS3Storer(
		cfg.BlobStorageBucketName, cfg.BlobStorageCreateBucket,
	)
	if err != nil {
		return nil, err
	}

	// setup http handlers
	handle := handler{storageClient: storageClient}

	// Register all routes
	// See routes.go
	app.Get(pathRoot, func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Get(pathBlobsDigest, errEndpointNotImplemented)
	app.Head(pathBlobsDigest, handle.BlobExists)
	app.Get(pathManifestsReference, errEndpointNotImplemented)
	app.Post(pathBlobsUploads, handle.BlobUploadLocation)
	app.Patch(pathBlobsUploadsReference, handle.BlobUpload)
	app.Put(pathBlobsUploadsReference, handle.BlobRefClose)
	app.Put(pathManifestsReference, func(c fiber.Ctx) error {
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
		err = storageClient.WriteManifest(
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
	})
	app.Get(pathTagsList, errEndpointNotImplemented)
	app.Delete(pathManifestsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathReferrersDigest, errEndpointNotImplemented)
	app.Get(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsUploadsReference, errEndpointNotImplemented)

	return app, nil
}
