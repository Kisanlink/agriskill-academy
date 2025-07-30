package db

import (
	"context"
	"fmt"
	"io"
	"time"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3File represents a file in S3
type S3File struct {
	Key         string    `json:"key"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// Filter represents a filter for S3 operations
type Filter struct {
	Field    string
	Operator FilterOp
	Value    interface{}
}

// FilterOp represents filter operations
type FilterOp string

const (
	FilterOpEqual FilterOp = "eq"
)

// Config represents S3 configuration
type Config struct {
	S3Region          string
	S3Bucket          string
	S3Endpoint        string
	S3ForcePathStyle  bool
	S3DisableSSL      bool
	S3AccessKeyID     string
	S3SecretAccessKey string
	LogLevel          string
}

// S3Manager provides S3 operations
type S3Manager struct {
	config     *Config
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// NewS3Manager creates a new S3 manager
func NewS3Manager(config *Config, logger interface{}) *S3Manager {
	return &S3Manager{
		config: config,
	}
}

// Connect establishes connection to S3
func (s *S3Manager) Connect(ctx context.Context) error {
	// Create AWS session
	awsConfig := &aws.Config{
		Region: aws.String(s.config.S3Region),
		Credentials: credentials.NewStaticCredentials(
			s.config.S3AccessKeyID,
			s.config.S3SecretAccessKey,
			"",
		),
	}

	// Configure endpoint for MinIO or custom S3-compatible service
	if s.config.S3Endpoint != "" {
		awsConfig.Endpoint = aws.String(s.config.S3Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(s.config.S3ForcePathStyle)
		awsConfig.DisableSSL = aws.Bool(s.config.S3DisableSSL)
	}

	// Create session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create S3 client
	s.s3Client = s3.New(sess)

	// Create uploader
	s.uploader = s3manager.NewUploader(sess)

	// Create downloader
	s.downloader = s3manager.NewDownloader(sess)

	// Test connection by listing buckets
	_, err = s.s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to connect to S3: %w", err)
	}

	return nil
}

// UploadFile uploads a file to S3
func (s *S3Manager) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string, metadata map[string]string) error {
	// Prepare upload input
	input := &s3manager.UploadInput{
		Bucket:      aws.String(s.config.S3Bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	// Add metadata if provided
	if metadata != nil {
		input.Metadata = aws.StringMap(metadata)
	}

	// Upload file
	_, err := s.uploader.UploadWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from S3
func (s *S3Manager) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	// Create a buffer to store the downloaded file
	buf := aws.NewWriteAtBuffer([]byte{})

	// Download file to buffer
	_, err := s.downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	// Return a reader from the buffer
	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// Delete deletes a file from S3
func (s *S3Manager) Delete(ctx context.Context, key string) error {
	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

// List lists files in S3 with filters
func (s *S3Manager) List(ctx context.Context, filters []Filter, result interface{}) error {
	// Build list input
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.S3Bucket),
	}

	// Apply filters
	for _, filter := range filters {
		if filter.Field == "prefix" && filter.Operator == FilterOpEqual {
			input.Prefix = aws.String(filter.Value.(string))
		}
	}

	// List objects
	resp, err := s.s3Client.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to list objects from S3: %w", err)
	}

	// Convert to S3File slice
	files := make([]S3File, 0, len(resp.Contents))
	for _, obj := range resp.Contents {
		files = append(files, S3File{
			Key:         *obj.Key,
			Size:        *obj.Size,
			ContentType: "", // S3 doesn't return content type in list
			CreatedAt:   *obj.LastModified,
		})
	}

	// Set result
	if resultSlice, ok := result.(*[]S3File); ok {
		*resultSlice = files
	}

	return nil
}

// GetByKey gets a file by key
func (s *S3Manager) GetByKey(ctx context.Context, key string, result interface{}) error {
	// Get object head to check if it exists
	_, err := s.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Get object metadata
	resp, err := s.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %w", err)
	}

	// Create S3File
	file := S3File{
		Key:         key,
		Size:        *resp.ContentLength,
		ContentType: *resp.ContentType,
		CreatedAt:   time.Now(), // S3 doesn't return creation time in GetObject
	}

	// Set result
	if resultFile, ok := result.(*S3File); ok {
		*resultFile = file
	}

	return nil
}

// BuildFilter creates a filter for S3 operations
func (s *S3Manager) BuildFilter(field string, op FilterOp, value interface{}) Filter {
	return Filter{
		Field:    field,
		Operator: op,
		Value:    value,
	}
}
