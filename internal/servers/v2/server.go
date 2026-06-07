package serverv2

import (
	"github.com/gofiber/fiber/v3"
	ociv2 "github.com/jbarzegar/oci/internal/oci/v2"
)

func handlePullManifest(c fiber.Ctx) error {
	invalid := []Error{}

	// validate each item then append any potential invalid
	// params to a slice. Allows for fixing all potential issues on first pass
	// rather than failing-fast each time something happens
	name := c.Params("name", "")
	nameValid := ociv2.ValidateManifestName(name)
	if !nameValid {
		invalid = append(invalid,
			serverError(ERR_NAME_INVALID, "param name invalid", ""))
	}

	ref := c.Params("reference", "")
	refValid := ociv2.ValidateReference(ref)
	if !refValid {
		invalid = append(invalid,
			serverError(ERR_MANIFEST_INVALID, "invalid manifest", ""))
	}

	if len(invalid) > 0 {
		return handleErrorResponse(c, 401, invalid...)
	}

	return nil
}

func New() *fiber.App {
	app := fiber.New()

	// Register all routes
	// See routes.go
	app.Get(pathRoot, errEndpointNotImplemented)
	app.Get(pathBlobsDigest, errEndpointNotImplemented)
	app.Get(pathManifestsReference, errEndpointNotImplemented)
	app.Post(pathBlobsUploads, errEndpointNotImplemented)
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
