package serverv2

import (
	"github.com/codingconcepts/env"
	"github.com/gofiber/fiber/v3"
	"github.com/jbarzegar/oci/internal/storagedriver"
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

	// init s3 driver and various deps
	storageClient, manifestWriter, manifestReader, err := storagedriver.
		InitS3Driver(
			cfg.BlobStorageBucketName,
			cfg.BlobStorageCreateBucket,
		)
	if err != nil {
		return nil, err
	}

	// setup http handlers
	handle := handler{
		storageClient: storageClient,
		manifest: manifestHandler{
			Reader: manifestReader,
			Writer: manifestWriter,
		},
	}

	app := fiber.New()
	// Register all routes
	// See routes.go
	app.Get(pathRoot, func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Get(pathBlobsDigest, errEndpointNotImplemented)
	app.Head(pathBlobsDigest, handle.BlobExists)
	app.Get(pathManifestsReference, handle.ManifestExists)
	app.Post(pathBlobsUploads, handle.BlobUploadLocation)
	app.Patch(pathBlobsUploadsReference, handle.BlobUpload)
	app.Put(pathBlobsUploadsReference, handle.BlobRefClose)
	app.Put(pathManifestsReference, handle.UploadManifest)
	app.Get(pathTagsList, errEndpointNotImplemented)
	app.Delete(pathManifestsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathReferrersDigest, errEndpointNotImplemented)
	app.Get(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsUploadsReference, errEndpointNotImplemented)

	return app, nil
}
