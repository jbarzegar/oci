package serverv2

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

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
	// Not adherant to spec
	ERR_NOT_IMPLEMENTED ServerErrCode = "NOT_IMPLEMENTED"
)

type Error struct {
	Code    ServerErrCode `json:"code"`
	Message string        `json:"message"`
	Detail  any           `json:"detail"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// serverError formats a server error based on a known list of
// error codes serverError should adhere to the oci distribution
// spec of errors and their shape
func serverError(
	code ServerErrCode,
	message string,
	detail any,
) Error {
	return Error{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// handleErrorResponse sends a response multiple server errors in a oci v2
// compliant fashion. Marshaling JSON in accordance with the spec (if requester accepts JSON)
// and sending it with a provided http status
func handleErrorResponse(ctx fiber.Ctx, status int, errs ...Error) error {
	if ctx.AcceptsJSON() {
		e := ErrorResponse{Errors: errs}
		payload, err := json.Marshal(e)
		if err != nil {
			return err
		}
		return ctx.Status(status).Send(payload)
	}

	// map errors to strings
	// CODE::Message
	// -- <detail>
	errStrings := []string{}
	for _, e := range errs {
		errStrings = append(errStrings, fmt.Sprintf("%v::%v\n-- %v", e.Code, e.Message, e.Detail))
	}

	return ctx.Status(status).Send([]byte(strings.Join(errStrings, "\n")))
}

// errEndpointNotImplemented is a catch all err func to
// communicate that a given endpoint hasn't yet been implemented
// this doesn't follow the spec but shows active lack-of
// compliance in a clear fashion
func errEndpointNotImplemented(c fiber.Ctx) error {
	return handleErrorResponse(c, fiber.ErrNotImplemented.Code,
		serverError(ERR_NOT_IMPLEMENTED,
			"The given route has not been implemented yet",
			nil,
		),
	)

}
