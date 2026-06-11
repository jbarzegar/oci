package blobstorage

import (
	"context"
	"io"
)

// Storer is a interface that simplifes what needs to be passed
// to a blob storage client currently doesn't account for access
// control
type Storer interface {
	GetImage(ctx context.Context, name string, digest string) (any, error)
	// StageImage stages the location an image may be saved to in the future
	StageImage(ctx context.Context, name string, digest string) (any, error)
	// Write image to given id
	WriteImage(ctx context.Context, name string, digest string, w io.Reader) error
	HasImage(ctx context.Context, name string, digest string) (bool, error)
}
