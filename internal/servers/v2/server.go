package serverv2

import (
	"github.com/codingconcepts/env"
	"github.com/gofiber/fiber/v3"
	"github.com/jbarzegar/oci/internal/blobstorage"
)

type V2ServerConfig struct {
	BlobStorageBucketName   string `env:"BLOB_STORAGE_BUCKET_NAME" required:"true"`
	BlobStorageCreateBucket bool   `env:"BLOB_STORAGE_CREATE_BUCKET" default:"false"`
}

func loadServerConfig() (*V2ServerConfig, error) {
	serverConfig := V2ServerConfig{}
	err := env.Set(&serverConfig)

	if err != nil {
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
	app.Get(pathRoot, errEndpointNotImplemented)
	app.Get(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathManifestsReference, errEndpointNotImplemented)
	app.Post(pathBlobsUploads, handle.BlobUploads)
	app.Patch(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Put(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Put(pathManifestsReference, errEndpointNotImplemented)
	app.Get(pathTagsList, errEndpointNotImplemented)
	app.Delete(pathManifestsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathReferrersDigest, errEndpointNotImplemented)
	app.Get(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsUploadsReference, errEndpointNotImplemented)

	return app, nil
}
