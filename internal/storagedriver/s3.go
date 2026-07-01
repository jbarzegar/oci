package storagedriver

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/codingconcepts/env"
	"github.com/jbarzegar/oci/internal/blobstorage"
	"github.com/jbarzegar/oci/internal/manifest"
)

type S3ClientConfig struct {
	Region          string `env:"RUSTFS_REGION" required:"true"`
	AccessKeyId     string `env:"RUSTFS_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string `env:"RUSTFS_SECRET_ACCESS_KEY" required:"true"`
	Endpoint        string `env:"RUSTFS_ENDPOINT_URL" required:"true"`
}

// loadS3ClientConfig will load the client config from the
// S3ClientConfig struct
func loadS3ClientConfig() (*S3ClientConfig, error) {
	config := S3ClientConfig{}

	if err := env.Set(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// initS3Client intializes and returns a aws client that enables
// s3 functionality from a given client config. Or returns an
// error if the client couldn't be initialized
func initS3Client(cfg *S3ClientConfig) (*s3.Client, error) {
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

func InitS3Driver(bucketName string, createBucket bool) (
	blobstorage.Storer,
	manifest.Writer,
	manifest.Reader,
	error,
) {
	s3Cfg, err := loadS3ClientConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	s3Client, err := initS3Client(s3Cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	storageClient, err := blobstorage.NewS3Storer(
		bucketName,
		s3Client,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	manifestReader, err := manifest.NewS3Reader(bucketName, s3Client)
	if err != nil {
		return nil, nil, nil, err
	}

	manifestWriter, err := manifest.NewS3Writer(
		bucketName,
		s3Client,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if createBucket {
		input := s3.CreateBucketInput{
			Bucket: &bucketName,
		}
		if _, err := s3Client.CreateBucket(
			context.Background(),
			&input,
		); err != nil {
			return nil, nil, nil, err
		}
	}

	return storageClient, manifestWriter, manifestReader, nil
}
