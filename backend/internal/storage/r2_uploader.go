// File: r2_uploader.go
// Purpose: Implement Cloudflare R2 uploader and build uploader instances from integration JSON config.
// Module: backend/internal/storage, object storage infrastructure layer.
// Related: app bootstrap uploader wiring, admin upload handler, integration provider configuration.
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Config struct {
	AccountID       string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	PublicBaseURL   string
}

type r2Uploader struct {
	client        *s3.Client
	bucket        string
	publicBaseURL string
	accountID     string
}

func NewR2UploaderFromConfigJSON(configJSON []byte) (Uploader, error) {
	cfg, err := parseR2Config(configJSON)
	if err != nil {
		return nil, err
	}
	return NewR2Uploader(cfg)
}

func NewR2Uploader(cfg R2Config) (Uploader, error) {
	if strings.TrimSpace(cfg.AccountID) == "" || strings.TrimSpace(cfg.Bucket) == "" || strings.TrimSpace(cfg.AccessKeyID) == "" || strings.TrimSpace(cfg.SecretAccessKey) == "" {
		return nil, fmt.Errorf("invalid r2 config: missing required fields")
	}

	endpointURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
		awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{URL: endpointURL, HostnameImmutable: true}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	return &r2Uploader{
		client:        client,
		bucket:        cfg.Bucket,
		publicBaseURL: strings.TrimRight(cfg.PublicBaseURL, "/"),
		accountID:     cfg.AccountID,
	}, nil
}

func (u *r2Uploader) Upload(ctx context.Context, objectKey string, body io.Reader, contentType string) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(objectKey),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("r2 put object failed: %w", err)
	}

	if u.publicBaseURL != "" {
		return u.publicBaseURL + "/" + strings.TrimLeft(objectKey, "/"), nil
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", u.accountID, u.bucket, strings.TrimLeft(objectKey, "/")), nil
}

func parseR2Config(configJSON []byte) (R2Config, error) {
	var payload map[string]any
	if err := json.Unmarshal(configJSON, &payload); err != nil {
		return R2Config{}, fmt.Errorf("invalid r2 config json: %w", err)
	}
	cfg := R2Config{
		AccountID:       strings.TrimSpace(getString(payload, "account_id")),
		Bucket:          strings.TrimSpace(getString(payload, "bucket")),
		AccessKeyID:     strings.TrimSpace(getString(payload, "access_key_id")),
		SecretAccessKey: strings.TrimSpace(getString(payload, "secret_access_key")),
		PublicBaseURL:   strings.TrimSpace(getString(payload, "public_base_url")),
	}
	if cfg.AccountID == "" || cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return R2Config{}, fmt.Errorf("invalid r2 config: account_id, bucket, access_key_id, secret_access_key are required")
	}
	return cfg, nil
}

func getString(payload map[string]any, key string) string {
	value, ok := payload[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
