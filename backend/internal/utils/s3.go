package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"ithozyeva/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type contextKey string

const contentSha256Key contextKey = "contentSha256"

type S3Client struct {
	client *s3.Client
	bucket string
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

	return &S3Client{client: client, bucket: cfg.Bucket}, nil
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *S3Client) Upload(ctx context.Context, key string, content []byte, contentType string) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("S3 Upload error: bucket=%s, key=%s, error=%v", c.bucket, key, err)
		return fmt.Errorf("failed to upload to S3: %w", err)
	}
	return nil
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
