package blobstorage

import (
	"bufio"
	"bytes"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Store existing writers in memory Not sure if this going to be
// the best way to store this info We'll do it for now. But I'd
// like to know if there's a more effective way of doing this
var writers = map[string]Writer{}
var writerMutex = sync.RWMutex{}

// NewS3Writer generates a new writer and stores it in a map if
// the writer exists by name already. the respective writer is
// returned instead
func NewS3Writer(
	name string,
	bucket string,
	uploadID uuid.UUID,
	s3Client *s3.Client,
) Writer {
	w, ok := writers[name]
	if !ok {
		go func() {
			writerMutex.Lock()
			w = &S3Writer{
				name:     name,
				uploads:  []uuid.UUID{uploadID},
				s3client: s3Client,
				bucket:   bucket,
				data:     []byte{},
			}
			writers[name] = w
			writerMutex.Unlock()
		}()
	}

	return w
}

type S3Writer struct {
	name     string
	bucket   string
	uploads  []uuid.UUID
	s3client *s3.Client
	data     []byte
}

// -- Implement S3 Writer --

func (w *S3Writer) AppendPart(
	ctx context.Context,
	p uuid.UUID,
	data []byte,
) (int, error) {
	w.data = append(w.data, data...)
	// calculate the length of bytes appended
	b := bufio.NewWriter(&bytes.Buffer{})
	nn, err := b.Write(w.data)
	if err != nil {
		return -1, err
	}

	return nn, nil
}

func (w *S3Writer) Write(ctx context.Context, digest string) error {
	r := bytes.NewBuffer(w.data)
	if r.Len() == 0 {
		return ErrNoWritableContent
	}
	key := formatKey(w.name, digest)
	input := s3.PutObjectInput{
		Bucket: &w.bucket,
		Key:    &key,
		Body:   bytes.NewReader(w.data),
	}
	_, err := w.s3client.PutObject(ctx, &input)
	return err
}

func (w *S3Writer) Flush() {
	w.data = make([]byte, 0)
}

// -- End Impl --
