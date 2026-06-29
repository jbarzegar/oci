package blobstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/codingconcepts/env"
	"github.com/google/uuid"
	"github.com/jbarzegar/oci/internal/manifest"
)

type s3ClientConfig struct {
	Region          string `env:"RUSTFS_REGION" required:"true"`
	AccessKeyId     string `env:"RUSTFS_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string `env:"RUSTFS_SECRET_ACCESS_KEY" required:"true"`
	Endpoint        string `env:"RUSTFS_ENDPOINT_URL" required:"true"`
}

// loadS3ClientConfig will load the client config from the
// S3ClientConfig struct
func loadS3ClientConfig() (*s3ClientConfig, error) {
	config := s3ClientConfig{}

	if err := env.Set(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// initS3Client intializes and returns a aws client that enables
// s3 functionality from a given client config. Or returns an
// error if the client couldn't be initialized
func initS3Client(cfg *s3ClientConfig) (*s3.Client, error) {
	credProvider := credentials.
		NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.SecretAccessKey, "")
	creds := aws.NewCredentialsCache(credProvider)
	resolver := aws.EndpointResolverFunc(
		func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: cfg.Endpoint,
			}, nil
		})
	// build aws.Config
	config := aws.Config{
		Region:           cfg.Region,
		EndpointResolver: resolver,
		Credentials:      creds,
	}

	// build S3 client
	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return client, nil
}

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

func NewS3Storer(bucketName string, createBucket bool) (Storer, error) {
	cfg, err := loadS3ClientConfig()
	if err != nil {
		return nil, err
	}
	client, err := initS3Client(cfg)
	if err != nil {
		return nil, err
	}

	if createBucket {
		input := s3.CreateBucketInput{
			Bucket: &bucketName,
		}
		if _, err := client.CreateBucket(context.Background(), &input); err != nil {
			return nil, err
		}
	}

	storer := &s3Storer{client: client, bucketName: bucketName}
	return storer, nil
}

// -- Impl s3Storer --

func (s *s3Storer) CreateWriter(
	ctx context.Context,
	name string,
	ref uuid.UUID,
) (Writer, error) {
	defer ctx.Done()
	input := s3.CreateMultipartUploadInput{
		Key:         &name,
		Bucket:      &s.bucketName,
		ContentType: aws.String("application/octet-stream"),
	}
	mu, err := s.client.CreateMultipartUpload(ctx, &input)
	if err != nil {
		return nil, err
	}
	w := newS3Writer(name, s.bucketName, ref, s.client, *mu.UploadId)
	return w, nil
}

func (s *s3Storer) GetWriterByUUID(
	ctx context.Context,
	ref uuid.UUID,
) (Writer, error) {
	for _, w := range writers {
		if slices.Contains(w.Parts(), ref) {
			return w, nil
		}
	}

	return nil, ErrWriterNotFound

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

func formatManifestKey(name string, ref string) string {
	n := fmt.Sprintf("%v/manifests", name)
	return formatKey(n, ref)
}

type manifestPayload struct {
	Key       string
	MediaType string
	Body      []byte
}

func (s *s3Storer) generateManifestInput(p manifestPayload) *s3.PutObjectInput {
	input := s3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &p.Key,
		ContentType: aws.String(p.MediaType),
		Body:        bytes.NewReader(p.Body),
	}

	return &input
}

func (s *s3Storer) WriteManifest(
	ctx context.Context,
	name string,
	tag string,
	m manifest.ManifestV2,
	mediaType string,
) error {
	b, err := manifest.MarshalV2(m)
	if err != nil {
		return err
	}

	// 	prepare manifest payload
	payload := manifestPayload{
		MediaType: mediaType,
		Body:      b,
		// Initial key is set to digest here.
		// Later on the Key may mutate to a named reference
		Key: formatManifestKey(name, m.Config.Digest),
	}

	// save raw digest.
	// A named reference _may_ change since tags are implicitly
	// mutable.
	// While tags / references are mutable the digests should be
	// saved independently. For preservation and or rollback, re-tags, etc
	if _, err = s.client.PutObject(
		ctx,
		s.generateManifestInput(payload),
	); err != nil {
		return err
	}

	// save manifest based on a tag / reference
	// FIXME(?): It would be nice to not have to save duplicate
	// data in the future and instead have a pointer somewhere
	// but manifests are not generally, extremely large.
	payload.Key = formatManifestKey(name, tag)
	if _, err = s.client.PutObject(
		ctx,
		s.generateManifestInput(payload),
	); err != nil {
		return err
	}

	return nil
}

func (s *s3Storer) ManifestInfo(ctx context.Context,
	name string,
	digest string,
) (*manifest.ManifestV2, *StoreInfo, bool, error) {
	key := formatManifestKey(name, digest)
	input := s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}

	output, err := s.client.GetObject(ctx, &input)
	if err != nil {
		return nil, nil, false, err
	}

	buf, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, nil, false, err
	}

	m, err := manifest.UnmarshalV2(buf)
	if err != nil {
		return nil, nil, false, err
	}

	info := StoreInfo{
		ContentType:   output.ContentType,
		ContentLength: output.ContentLength,
	}

	return &m, &info, true, nil
}

// -- end impl --
