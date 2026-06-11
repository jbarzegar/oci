package blobstorage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/codingconcepts/env"
)

type S3ClientConfig struct {
	Region          string `env:"RUSTFS_REGION" required:"true"`
	AccessKeyId     string `env:"RUSTFS_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string `env:"RUSTFS_SECRET_ACCESS_KEY" required:"true"`
	Endpoint        string `env:"RUSTFS_ENDPOINT_URL" required:"true"`
}

func LoadConfig() (*S3ClientConfig, error) {
	config := S3ClientConfig{}

	if err := env.Set(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func initS3Client(cfg *S3ClientConfig) (*s3.Client, error) {
	// build aws.Config
	credProvider := credentials.
		NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.SecretAccessKey, "")
	creds := aws.NewCredentialsCache(credProvider)
	resolver := aws.EndpointResolverFunc(
		func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: cfg.Endpoint,
			}, nil
		})
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

func NewS3Storer(bucketName string, createBucket bool) (Storer, error) {
	cfg, err := LoadConfig()
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

type s3Storer struct {
	bucketName string
	client     *s3.Client
}

// formatKey formats a container name and diget into a directory
// format. Joining the two strings with a slash within the object
// store. This appears as seperate directories so an image will
// be collected similar to /<name> -- /<name>/<digest>
func formatKey(name string, digest string) string {
	return fmt.Sprintf("%v/%v", name, digest)
}

func (s *s3Storer) GetImage(
	ctx context.Context,
	name string, digest string,
) (any, error) {
	key := formatKey(name, digest)
	input := s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}
	// TODO: Gotta do something with the blob idk what to do with
	// it
	_, err := s.client.GetObject(ctx, &input)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *s3Storer) StageImage(
	ctx context.Context,
	name string, digest string,
) (any, error) {
	return nil, nil
}

func (s *s3Storer) WriteImage(
	ctx context.Context,
	name string, digest string,
	w io.Reader,
) error {
	key := formatKey(name, digest)
	input := s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
		Body:   w,
	}
	_, err := s.client.PutObject(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}

func (s *s3Storer) HasImage(
	ctx context.Context,
	name string, digest string,
) (bool, error) {
	return false, nil
}
