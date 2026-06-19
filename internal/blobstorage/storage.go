package blobstorage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jbarzegar/oci/internal/manifest"
)

// Storer is a interface that simplifes what needs to be passed
// to a blob storage client currently doesn't account for access
// control
type Storer interface {
	// Create an instance of Writer based on the name of an image
	CreateWriter(ctx context.Context, name string, ref uuid.UUID) (Writer, error)
	// Fetch a given writer by a reference uuid this is a O(n)
	// operation Since the operation requires finding by a given
	// session ID Otherwise will error with a unfound writer
	// error
	GetWriterByUUID(ctx context.Context, ref uuid.UUID) (Writer, error)
	// GetWriterByName fetches a given writer This method is a
	// O(1) operation as finding the writer requires registering
	GetWriterByName(ctx context.Context, name string) (Writer,
		error)
	// BlobInfo checks if a given name & digest has been written
	// returns size of the blob in bytes, whether it exists, and
	// potentiall error
	BlobInfo(ctx context.Context, name string, digest string) (int64, bool, error)
	// WriteManfest writes a manifest to a given location with a digest
	WriteManifest(ctx context.Context, name string, digest string, manifest manifest.ManifestV2) error
}

var ErrWriterNotFound = errors.New("writer not found")

// Writer is a interface that controls how blobs are writen while handled
type Writer interface {
	// Name returns the repo name the Writer was registered with
	Name() string
	// AppendPart adds a new upload to the parts list
	// Returns the content length (in bytes) appened
	AppendPart(ctx context.Context, p uuid.UUID, data []byte) (int, error)
	// Write writes all parts to a given source
	Write(ctx context.Context, digest string) error
	// Parts returns the registered uploads for a given writer instance
	Parts() []uuid.UUID
	// UploadID returns the id of the given writers upload instance
	UploadID() *string
	// Flush will purge all parts of a given writer instance this
	// should be done to cleanup after a writer does a final
	// write to its store
	Flush()
}
