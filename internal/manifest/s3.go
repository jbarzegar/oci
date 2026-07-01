package manifest

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Reader struct {
	bucketName string
	client     *s3.Client
}

func NewS3Reader(bucketName string, client *s3.Client) (Reader, error) {
	return &S3Reader{bucketName: bucketName, client: client}, nil
}

func (s *S3Reader) Get(
	ctx context.Context,
	input *ReaderGetInput,
) (*ReaderGetOutput, bool, error) {
	key := formatManifestKey(input.Name, input.Digest)

	s3Input := s3.GetObjectInput{
		Key:    &key,
		Bucket: &s.bucketName,
	}
	output, err := s.client.GetObject(ctx, &s3Input)
	if err != nil {
		return nil, false, err
	}

	info := Content{
		Type:   output.ContentType,
		Length: output.ContentLength,
	}
	return &ReaderGetOutput{
		Body:    output.Body,
		Content: &info,
	}, true, nil
}

// --
type manifestPayload struct {
	Key       string
	MediaType string
	Body      []byte
}

func (w *S3Writer) generateManifestInput(p manifestPayload) *s3.PutObjectInput {
	input := s3.PutObjectInput{
		Bucket:      &w.bucketName,
		Key:         &p.Key,
		ContentType: aws.String(p.MediaType),
		Body:        bytes.NewReader(p.Body),
	}

	return &input
}

func NewS3Writer(bucketName string, client *s3.Client) (Writer, error) {
	return &S3Writer{bucketName: bucketName, client: client}, nil
}

type S3Writer struct {
	bucketName string
	client     *s3.Client
}

func (w *S3Writer) Write(
	ctx context.Context,
	input *WriterWriteInput,
) error {
	b, err := MarshalV2(input.Manifest)
	if err != nil {
		return err
	}

	// 	prepare manifest paylod
	payload := manifestPayload{
		MediaType: input.MediaType,
		Body:      b,
		// Initial key is set to digest here.
		// Later on the Key may mutate to a named reference
		Key: formatManifestKey(
			input.Name,
			input.Manifest.Config.Digest,
		),
	}

	// save raw digest.
	// A named reference _may_ change since tags are implicitly
	// mutable.
	// While tags / references are mutable the digests should be
	// saved independently. For preservation and or rollback, re-tags, etc
	if _, err = w.client.PutObject(
		ctx,
		w.generateManifestInput(payload),
	); err != nil {
		return err
	}

	// save manifest based on a tag / reference
	// FIXME(?): It would be nice to not have to save duplicate
	// data in the future and instead have a pointer somewhere
	// but manifests are not generally, extremely large.
	payload.Key = formatManifestKey(input.Name, input.Tag)
	if _, err = w.client.PutObject(
		ctx,
		w.generateManifestInput(payload),
	); err != nil {
		return err
	}

	return nil
}
