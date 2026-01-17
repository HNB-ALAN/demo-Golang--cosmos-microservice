package database

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
)

// MinIOClient interface for MinIO operations
type MinIOClient interface {
	PutObject(ctx context.Context, bucketName, objectName string, data []byte) error
	GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
	DeleteObject(ctx context.Context, bucketName, objectName string) error
	ListObjects(ctx context.Context, bucketName, prefix string) ([]string, error)
	Ping(ctx context.Context) error
}

// minioClient implements MinIOClient interface
type minioClient struct {
	client *minio.Client
	logger logging.Logger
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(cfg *config.Config, logger logging.Logger) (MinIOClient, error) {
	if !cfg.MinIO.Enabled {
		return nil, fmt.Errorf("minio is not enabled in configuration")
	}

	// Create MinIO client with host:port
	endpoint := fmt.Sprintf("%s:%d", cfg.MinIO.Host, cfg.MinIO.Port)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.SSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &minioClient{
		client: client,
		logger: logger,
	}, nil
}

// PutObject uploads an object to MinIO
func (c *minioClient) PutObject(ctx context.Context, bucketName, objectName string, data []byte) error {
	reader := strings.NewReader(string(data))
	_, err := c.client.PutObject(ctx, bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		c.logger.Error("MinIO PutObject failed",
			logging.String("bucket", bucketName),
			logging.String("object", objectName),
			logging.Error(err))
		return fmt.Errorf("minio put object failed: %w", err)
	}
	return nil
}

// GetObject downloads an object from MinIO
func (c *minioClient) GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	object, err := c.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.logger.Error("MinIO GetObject failed",
			logging.String("bucket", bucketName),
			logging.String("object", objectName),
			logging.Error(err))
		return nil, fmt.Errorf("minio get object failed: %w", err)
	}
	defer object.Close()

	// Read all data
	buf := make([]byte, 0, 1024*1024) // Start with 1MB buffer
	for {
		n, err := object.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if err != nil {
			if err == io.EOF {
				break
			}
			c.logger.Error("MinIO Read failed",
				logging.String("bucket", bucketName),
				logging.String("object", objectName),
				logging.Error(err))
			return nil, fmt.Errorf("minio read failed: %w", err)
		}
		if len(buf) == cap(buf) {
			// Double buffer size
			newBuf := make([]byte, len(buf), 2*cap(buf))
			copy(newBuf, buf)
			buf = newBuf
		}
	}

	return buf, nil
}

// DeleteObject deletes an object from MinIO
func (c *minioClient) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	err := c.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.logger.Error("MinIO DeleteObject failed",
			logging.String("bucket", bucketName),
			logging.String("object", objectName),
			logging.Error(err))
		return fmt.Errorf("minio delete object failed: %w", err)
	}
	return nil
}

// ListObjects lists objects in a bucket with prefix
func (c *minioClient) ListObjects(ctx context.Context, bucketName, prefix string) ([]string, error) {
	objectCh := c.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var objects []string
	for object := range objectCh {
		if object.Err != nil {
			c.logger.Error("MinIO ListObjects failed",
				logging.String("bucket", bucketName),
				logging.String("prefix", prefix),
				logging.Error(object.Err))
			return nil, fmt.Errorf("minio list objects failed: %w", object.Err)
		}
		objects = append(objects, object.Key)
	}

	return objects, nil
}

// Ping checks MinIO connection
func (c *minioClient) Ping(ctx context.Context) error {
	// Try to list buckets to check connection
	_, err := c.client.ListBuckets(ctx)
	if err != nil {
		c.logger.Error("MinIO Ping failed", logging.Error(err))
		return fmt.Errorf("minio ping failed: %w", err)
	}
	return nil
}

// MinIOHealthChecker implements HealthChecker for MinIO
type MinIOHealthChecker struct {
	client MinIOClient
	logger logging.Logger
}

// NewMinIOHealthChecker creates a new MinIO health checker
func NewMinIOHealthChecker(client MinIOClient, logger logging.Logger) *MinIOHealthChecker {
	return &MinIOHealthChecker{
		client: client,
		logger: logger,
	}
}

// Check performs health check for MinIO
func (h *MinIOHealthChecker) Check(ctx context.Context) error {
	// Set timeout for health check
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.client.Ping(ctx); err != nil {
		h.logger.Error("MinIO health check failed", logging.Error(err))
		return fmt.Errorf("minio health check failed: %w", err)
	}

	return nil
}

// initializeMinIO initializes MinIO connection with retry logic
func (m *DatabaseManager) initializeMinIO() error {
	maxRetries := 5
	baseDelay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := baseDelay * time.Duration(1<<attempt)
			time.Sleep(delay)
		}

		// Create no-op logger for MinIO client
		emptyLogger := logging.NewLogger("minio", config.LogConfig{})
		client, err := NewMinIOClient(m.config, *emptyLogger)
		if err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("failed to create MinIO client after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Ping(ctx); err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("failed to connect to MinIO after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		m.minio = client
		return nil
	}

	return fmt.Errorf("failed to initialize MinIO after %d attempts", maxRetries)
}
