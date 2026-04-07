package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"ithozyeva/config"
	"ithozyeva/internal/s3resolve"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type S3Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewS3Client() (*S3Client, error) {
	cfg := config.CFG.S3
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("s3 config is not fully specified")
	}

	// Проверяем и нормализуем endpoint
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("S3 endpoint is required")
	}

	// Парсим URL для проверки корректности
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid S3 endpoint URL: %w", err)
	}

	// Если протокол не указан, добавляем https
	if parsedURL.Scheme == "" {
		endpoint = "https://" + endpoint
		parsedURL, err = url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid S3 endpoint URL after adding scheme: %w", err)
		}
	}

	endpointBase := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = "ru-central-1"
	}

	accessKey := strings.TrimSpace(cfg.AccessKey)
	secretKey := strings.TrimSpace(cfg.SecretKey)

	awsCfg := aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpointBase)
		o.UsePathStyle = true
		o.UseARNRegion = false
		o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
			if err := stack.Build.Add(cleanupHeadersMiddleware{}, middleware.After); err != nil {
				return err
			}
			// Логируем запрос ПОСЛЕ формирования подписи (в Finalize)
			return stack.Finalize.Add(loggingMiddleware{}, middleware.After)
		})
	})

	presignClient := s3.NewPresignClient(client)

	return &S3Client{
		client:        client,
		presignClient: presignClient,
		bucket:        cfg.Bucket,
	}, nil
}

// cleanupHeadersMiddleware удаляет лишние заголовки AWS SDK, которые могут мешать cloud.ru
type cleanupHeadersMiddleware struct{}

func (m cleanupHeadersMiddleware) ID() string {
	return "S3CleanupHeaders"
}

func (m cleanupHeadersMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	return next.HandleInitialize(ctx, in)
}

func (m cleanupHeadersMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	return next.HandleSerialize(ctx, in)
}

func (m cleanupHeadersMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		req.Header.Del("Amz-Sdk-Request")
		req.Header.Del("Amz-Sdk-Invocation-Id")
		req.Header.Del("User-Agent")
	}
	return next.HandleBuild(ctx, in)
}

func (m cleanupHeadersMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	return next.HandleFinalize(ctx, in)
}

// loggingMiddleware логирует HTTP запросы к S3 для отладки
type loggingMiddleware struct{}

func (m loggingMiddleware) ID() string {
	return "S3RequestLogging"
}

func (m loggingMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	return next.HandleInitialize(ctx, in)
}

func (m loggingMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	return next.HandleSerialize(ctx, in)
}

func (m loggingMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	return next.HandleBuild(ctx, in)
}

func (m loggingMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		var headerKeys []string
		canonicalHeaders := make(map[string]string)

		host := req.URL.Host
		if host == "" {
			host = req.Host
		}
		if host != "" {
			headerKeys = append(headerKeys, "host")
			canonicalHeaders["host"] = host
		}

		if req.ContentLength > 0 {
			headerKeys = append(headerKeys, "content-length")
			canonicalHeaders["content-length"] = fmt.Sprintf("%d", req.ContentLength)
		}

		for key, values := range req.Header {
			lowerKey := strings.ToLower(key)
			if lowerKey == "host" || lowerKey == "content-length" {
				continue
			}
			if _, exists := canonicalHeaders[lowerKey]; !exists {
				headerKeys = append(headerKeys, lowerKey)
			}
			canonicalHeaders[lowerKey] = strings.TrimSpace(strings.Join(values, ","))
		}

		for i := 0; i < len(headerKeys)-1; i++ {
			for j := i + 1; j < len(headerKeys); j++ {
				if headerKeys[i] > headerKeys[j] {
					headerKeys[i], headerKeys[j] = headerKeys[j], headerKeys[i]
				}
			}
		}
	}

	return next.HandleFinalize(ctx, in)
}

func (c *S3Client) Upload(ctx context.Context, key string, content []byte, contentType string) error {
	return c.UploadWithACL(ctx, key, content, contentType, "")
}

func (c *S3Client) UploadPublic(ctx context.Context, key string, content []byte, contentType string) error {
	return c.UploadWithACL(ctx, key, content, contentType, "public-read")
}

func (c *S3Client) UploadWithACL(ctx context.Context, key string, content []byte, contentType string, acl string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	}

	// ACL не обязателен для Cloud.ru, так как используются presigned URLs
	// Но оставляем для совместимости с AWS S3
	if acl != "" {
		input.ACL = s3types.ObjectCannedACL(acl)
	}

	_, err := c.client.PutObject(ctx, input)
	if err != nil {
		log.Printf("S3 Upload error: bucket=%s, key=%s, error=%v", c.bucket, key, err)
		return fmt.Errorf("failed to upload to S3: %w", err)
	}
	log.Printf("S3 Upload successful: bucket=%s, key=%s", c.bucket, key)
	return nil
}

// GetPublicURL returns just the S3 key for storage in the database.
func (c *S3Client) GetPublicURL(key string) string {
	return key
}

// GetPresignedURL generates a presigned URL for the given S3 key with 7-day expiry.
func (c *S3Client) GetPresignedURL(key string) (string, error) {
	return c.GetPresignedURLWithExpiry(key, 7*24*time.Hour)
}

// GetPresignedURLWithExpiry generates a presigned URL with custom expiry time.
func (c *S3Client) GetPresignedURLWithExpiry(key string, expiry time.Duration) (string, error) {
	result, err := c.presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}
	return result.URL, nil
}

// ResolveURL takes a stored value (S3 key or old-style full URL) and returns a presigned URL.
// Returns empty string if the input is empty.
func (c *S3Client) ResolveURL(stored string) string {
	if stored == "" {
		return ""
	}
	key := ExtractS3Key(stored, c.bucket)

	// Use presigned URL for all files (Cloud.ru requires authorization even for "public" files)
	url, err := c.GetPresignedURL(key)
	if err != nil {
		log.Printf("Failed to generate presigned URL for key %s: %v", key, err)
		return stored
	}
	return url
}

// ExtractS3Key extracts the S3 object key from either a full URL or a bare key.
func ExtractS3Key(stored string, bucket string) string {
	if !strings.HasPrefix(stored, "http") {
		return stored
	}
	parsed, err := url.Parse(stored)
	if err != nil {
		return stored
	}
	path := strings.TrimPrefix(parsed.Path, "/")
	path = strings.TrimPrefix(path, bucket+"/")
	return path
}

// Global S3 resolver for use in model hooks
var globalS3 *S3Client

func InitGlobalS3() {
	client, err := NewS3Client()
	if err != nil {
		log.Printf("Warning: failed to init global S3 client: %v", err)
		return
	}
	globalS3 = client
	// Register resolver for use in model hooks (breaks import cycle)
	s3resolve.Resolver = func(stored string) string {
		return client.ResolveURL(stored)
	}
}

// ResolveS3URL resolves a stored S3 key or old-style URL to a presigned URL.
func ResolveS3URL(stored string) string {
	if globalS3 == nil || stored == "" {
		return stored
	}
	return globalS3.ResolveURL(stored)
}

func (c *S3Client) Delete(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (c *S3Client) Download(ctx context.Context, key string) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, obj.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
