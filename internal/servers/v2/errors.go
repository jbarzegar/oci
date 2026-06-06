package serverv2

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

type errorCode string

const (
	ErrManifestBlobUnknown errorCode = "MANIFEST_BLOB_UNKNOWN"
)

// func formatError(code errorCode, message string, detail any) error {
// 	return nil
// }

func errEndpointNotImplemented(c fiber.Ctx) error {
	return c.
		Status(fiber.ErrNotImplemented.Code).
		SendString("Endpoint not yet implemented")

}

type ServerErrCode string

const (
	ERR_BLOB_UNKNOWN          ServerErrCode = "BLOB_UNKNOWN"
	ERR_BLOB_UPLOAD_INVALID   ServerErrCode = "BLOB_UPLOAD_INVALID"
	ERR_BLOB_UPLOAD_UNKNOWN   ServerErrCode = "BLOB_UPLOAD_UNKNOWN"
	ERR_DIGEST_INVALID        ServerErrCode = "DIGEST_INVALID"
	ERR_MANIFEST_BLOB_UNKNOWN ServerErrCode = "MANIFEST_BLOB_UNKNOWN"
	ERR_MANIFEST_INVALID      ServerErrCode = "MANIFEST_INVALID"
	ERR_MANIFEST_UNKNOWN      ServerErrCode = "MANIFEST_UNKNOWN"
	ERR_NAME_INVALID          ServerErrCode = "NAME_INVALID"
	ERR_NAME_UNKNOWN          ServerErrCode = "NAME_UNKNOWN"
	ERR_SIZE_INVALID          ServerErrCode = "SIZE_INVALID"
	ERR_UNAUTHORIZED          ServerErrCode = "UNAUTHORIZED"
	ERR_DENIED                ServerErrCode = "DENIED"
	ERR_UNSUPPORTED           ServerErrCode = "UNSUPPORTED"
	ERR_TOOMANYREQUESTS       ServerErrCode = "TOOMANYREQUESTS"
)

type ServerError struct {
	Code    ServerErrCode `json:"code"`
	Message string        `json:"message"`
	Detail  any           `json:"detail"`
}

type MultipleServerErrors struct {
	Errors []ServerError `json:"errors"`
}

// serverError formats a server error based on a known list of error codes
// serverError should adhere to the oci distribution spec of errors and their shape
func serverError(
	code ServerErrCode,
	message string,
	detail any,
) ServerError {
	return ServerError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// handleServerErrors handles multiple server errors in a oci v2 compliant fashion
// marshals JSON in accordance with the spec or throws an unhandled exception
func handleServerErrors(ctx fiber.Ctx, status int, errs ...ServerError) error {
	e := MultipleServerErrors{
		Errors: errs,
	}
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return ctx.Status(status).Send(payload)
}
