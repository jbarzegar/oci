package serverv2

import (
	"github.com/gofiber/fiber/v3"
)

func New() *fiber.App {
	app := fiber.New()

	// Register all routes
	// See routes.go
	app.Get(pathRoot, errEndpointNotImplemented)
	app.Get(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathManifestsReference, errEndpointNotImplemented)
	app.Post(pathBlobsUploads, handleBlobUploads)
	app.Patch(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Put(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Put(pathManifestsReference, errEndpointNotImplemented)
	app.Get(pathTagsList, errEndpointNotImplemented)
	app.Delete(pathManifestsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathReferrersDigest, errEndpointNotImplemented)
	app.Get(pathBlobsUploadsReference, errEndpointNotImplemented)
	app.Delete(pathBlobsUploadsReference, errEndpointNotImplemented)

	return app
}
