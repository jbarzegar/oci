package serverv2

// Houses all relevant routes for this version of server
// See below for full spec and methods
// https://github.com/opencontainers/distribution-spec/blob/main/spec.md#endpoints
const (
	pathRoot                  = "/v2/"
	pathBlobsDigest           = "/v2/:name/blobs/:digest"
	pathManifestsReference    = "/v2/:name/manifests/:reference"
	pathBlobsUploads          = "/v2/:name/blobs/uploads/"
	pathBlobsUploadsReference = "/v2/:name/blobs/uploads/:reference"
	pathTagsList              = "/v2/:name/tags/list"
	pathReferrersDigest       = "/v2/:name/referrers/:digest"
)
