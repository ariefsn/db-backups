package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client   *s3.Client
	bucketName string
}

type UploadMetadata struct {
	DatabaseType string
	Host         string
	Database     string
	Timestamp    time.Time
	FileSize     int64
}

// NewClient creates a new R2 storage client from environment variables
func NewClient() (*Client, error) {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	region := os.Getenv("R2_REGION")

	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" || bucketName == "" {
		return nil, fmt.Errorf("missing required R2 environment variables")
	}

	if region == "" {
		region = "auto"
	}

	// Create S3 client with custom endpoint for R2
	s3Client := s3.New(s3.Options{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		),
		BaseEndpoint: aws.String(endpoint),
	})

	return &Client{
		s3Client:   s3Client,
		bucketName: bucketName,
	}, nil
}

// Upload uploads a file to R2 and returns the object key
func (c *Client) Upload(ctx context.Context, filePath string, metadata UploadMetadata) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Generate object key: backups/{type}/{filename}
	filename := filepath.Base(filePath)
	objectKey := fmt.Sprintf("backups/%s/%s", metadata.DatabaseType, filename)

	// Prepare metadata for S3
	s3Metadata := map[string]string{
		"database-type": metadata.DatabaseType,
		"host":          metadata.Host,
		"database":      metadata.Database,
		"timestamp":     metadata.Timestamp.Format(time.RFC3339),
		"file-size":     fmt.Sprintf("%d", fileInfo.Size()),
	}

	// Upload to R2
	_, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:   aws.String(c.bucketName),
		Key:      aws.String(objectKey),
		Body:     file,
		Metadata: s3Metadata,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	return objectKey, nil
}

// Delete deletes an object from R2
func (c *Client) Delete(ctx context.Context, objectKey string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}

	return nil
}

// GetPresignedURL generates a presigned URL for downloading an object
func (c *Client) GetPresignedURL(ctx context.Context, objectKey string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}
