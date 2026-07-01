package blobstorage

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Storer is a interface that simplifes what needs to be passed
// to a blob storage client currently doesn't account for access
// control
type Storer interface {
	// Create an instance of Writer based on the name of an image
	CreateWriter(name string, ref uuid.UUID) (Writer, error)
	// GetWriterByName fetches a given writer This method is a
	// O(1) operation as finding the writer requires registering
	GetWriterByName(ctx context.Context, name string) (Writer,
		error)
	// BlobInfo checks if a given name & digest has been written
	// returns size of the blob in bytes, whether it exists, or
	// potential error. If an error is returned the preceeding
	// values should always be there 0
	BlobInfo(ctx context.Context, name string, digest string) (int64, bool, error)
}

var (
	ErrWriterNotFound    = errors.New("writer not found")
	ErrNoWritableContent = errors.New("no content to write")
)

// Writer is a interface that controls how blobs are writen while handled
type Writer interface {
	// AppendPart adds a new upload to the parts list
	// Returns the content length (in bytes) appened
	AppendPart(ctx context.Context, p uuid.UUID, data []byte) (int, error)
	// Write writes all parts to a given source
	Write(ctx context.Context, digest string) error
	// Flush will purge all parts of a given writer instance this
	// should be done to cleanup after a writer does a final
	// write to its store
	Flush()
}
