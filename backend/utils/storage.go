package utils

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// StorageService defines the contract for file storage operations
type StorageService interface {
	UploadFile(ctx context.Context, fileName string, fileReader io.Reader, contentType string) (string, error)
}

type ociStorageService struct {
	s3Client   *s3.Client
	bucketName string
	endpoint   string
	cdnURL     string
}

// NewOCIStorageService initializes a new OCI S3-compatible storage service
func NewOCIStorageService(accessKeyID, secretAccessKey, region, bucketName, endpoint, cdnURL string) (StorageService, error) {
	// Ensure endpoint starts with http:// or https:// to prevent empty protocol scheme error in AWS SDK
	if endpoint != "" && !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// If credentials are not set, fallback to a mockup storage service for easy local testing
	if accessKeyID == "" || secretAccessKey == "" {
		return &mockStorageService{}, nil
	}

	// S3 Resolver custom endpoint for OCI
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown service %s", service)
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // OCI Object Storage requires path-style routing
	})

	return &ociStorageService{
		s3Client:   s3Client,
		bucketName: bucketName,
		endpoint:   endpoint,
		cdnURL:     cdnURL,
	}, nil
}

// UploadFile uploads an image or file to the OCI Bucket and returns its URL
func (s *ociStorageService) UploadFile(ctx context.Context, fileName string, fileReader io.Reader, contentType string) (string, error) {
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	// Determine payload size to bypass AWS SDK chunked encoding for OCI compatibility
	var size int64
	if seeker, ok := fileReader.(io.Seeker); ok {
		currentPos, err := seeker.Seek(0, io.SeekCurrent)
		if err == nil {
			endPos, err := seeker.Seek(0, io.SeekEnd)
			if err == nil {
				size = endPos
				_, _ = seeker.Seek(currentPos, io.SeekStart) // Restore position
			}
		}
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(uniqueFileName),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	}
	if size > 0 {
		input.ContentLength = aws.Int64(size)
	}

	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload object to OCI bucket: %v", err)
	}

	baseURL := s.endpoint
	if s.cdnURL != "" {
		baseURL = s.cdnURL
	}
	// Ensure baseURL has protocol scheme
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	fileURL := fmt.Sprintf("%s/%s/%s", baseURL, s.bucketName, uniqueFileName)
	return fileURL, nil
}

type mockStorageService struct{}

func (m *mockStorageService) UploadFile(ctx context.Context, fileName string, fileReader io.Reader, contentType string) (string, error) {
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)
	return fmt.Sprintf("https://mock-oci-storage.com/tutora-bucket/%s", uniqueFileName), nil
}
