package s3_storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	client   *s3.Client
	bucket   string
	region   string
	endpoint string
}

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Endpoint        string
}

func NewS3Client(ctx context.Context, cfg Config) (*S3Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to load AWS SDK config: %w", err)
	}

	s3Opts := func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		}
	}

	return &S3Client{
		client:   s3.NewFromConfig(awsCfg, s3Opts),
		bucket:   cfg.Bucket,
		region:   cfg.Region,
		endpoint: cfg.Endpoint,
	}, nil
}

func (s *S3Client) Upload(ctx context.Context, folder string, filename string, file io.Reader, contentType string) (string, error) {
	objectKey := fmt.Sprintf("%s/%s", folder, filename)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to upload object to S3: %w", err)
	}

	var fileURL string
	if strings.Contains(s.endpoint, "supabase.co") {
		baseDomain := strings.Split(s.endpoint, "/storage")[0]
		fileURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", baseDomain, s.bucket, objectKey)
	} else {
		fileURL = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(s.endpoint, "/"), s.bucket, objectKey)
	}

	return fileURL, nil
}

func (s *S3Client) Delete(ctx context.Context, fileURL string) error {
	if fileURL == "" {
		return nil
	}

	var objectKey string

	if strings.Contains(fileURL, "supabase.co") {
		anchor := fmt.Sprintf("/public/%s/", s.bucket)
		parts := strings.Split(fileURL, anchor)
		if len(parts) < 2 {
			return fmt.Errorf("/invalid supabase file url format: anchor '%s' not found", anchor)
		}
		objectKey = parts[1]
	} else {
		expectedAnchor := fmt.Sprintf("/%s/", s.bucket)
		if !strings.Contains(fileURL, expectedAnchor) {
			expectedAnchor = fmt.Sprintf("%s.s3.", s.bucket)
		}

		parts := strings.Split(fileURL, expectedAnchor)
		if len(parts) < 2 {
			return fmt.Errorf("Invalid s3 file url format: bucket anchor not found")
		}
		objectKey = parts[1]
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("Failed to delete object with key = %s from S3: %w", objectKey, err)
	}

	return nil
}
