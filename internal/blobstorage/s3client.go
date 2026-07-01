package blobstorage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// formatKey formats a container name and diget into a directory
// format. Joining the two strings with a slash within the object
// store. This appears as seperate directories so an image will
// be collected similar to /<name> -- /<name>/<digest>
func formatKey(name string, digest string) string {
	return fmt.Sprintf("%v/%v", name, digest)
}

type s3Storer struct {
	bucketName string
	client     *s3.Client
}

func NewS3Storer(bucketName string, client *s3.Client) (Storer, error) {
	storer := &s3Storer{client: client, bucketName: bucketName}
	return storer, nil
}

// -- Impl s3Storer --

func (s *s3Storer) CreateWriter(
	name string,
	ref uuid.UUID,
) (Writer, error) {
	w := NewS3Writer(name, s.bucketName, ref, s.client)
	return w, nil
}

func (s *s3Storer) GetWriterByName(
	ctx context.Context,
	name string,
) (Writer, error) {
	defer ctx.Done()
	w, ok := writers[name]
	if !ok {
		return nil, ErrWriterNotFound
	}

	return w, nil
}

func (s *s3Storer) BlobInfo(
	ctx context.Context,
	name string,
	digest string,
) (int64, bool, error) {
	key := formatKey(name, digest)
	input := s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}
	x, err := s.client.GetObject(ctx, &input)
	if err != nil {
		return -1, false, err
	}

	return *x.ContentLength, true, nil
}

// -- end impl --
