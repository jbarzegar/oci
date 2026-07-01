package manifest

import (
	"context"
	"fmt"
	"io"
)

type Content struct {
	Type   *string
	Length *int64
}

type ReaderGetInput struct {
	Name   string
	Digest string
}

type ReaderGetOutput struct {
	Body    io.ReadCloser
	Content *Content
}

type Reader interface {
	// Get checks if a given name & digest has been
	// written returns manifest, whether it exists or a potential
	// error. If an error is returned the preceeding
	// values should always be there 0
	Get(ctx context.Context, input *ReaderGetInput) (*ReaderGetOutput, bool, error)
}

type WriterWriteInput struct {
	Name      string
	Tag       string
	Manifest  ManifestV2
	MediaType string
}

type Writer interface {
	// Write writes a manifest to a given location with a tag
	// Additionally the digest from the manifest will be written in turn.
	Write(ctx context.Context, input *WriterWriteInput) error
}

func formatManifestKey(name string, ref string) string {
	return fmt.Sprintf("%v/manifests/%v", name, ref)
}
