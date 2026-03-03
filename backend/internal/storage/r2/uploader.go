// File: uploader.go
// Purpose: Upload objects to Cloudflare R2 using S3-compatible AWS SDK.
// Module: backend/internal/storage/r2, storage provider layer.
// Related: config.R2Config and admin upload endpoints.
package r2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader struct {
	bucket        string
	publicBaseURL string
	client        *s3.Client
}

func New(ctx context.Context, cfg config.R2Config) (*Uploader, error) {
	if !cfg.Enabled {
		return nil, errors.New("r2 is disabled")
	}
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		return nil, errors.New("r2 config is incomplete")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)
	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = awsv2.String(endpoint)
		o.UsePathStyle = true
	})

	return &Uploader{bucket: cfg.Bucket, publicBaseURL: strings.TrimRight(cfg.PublicBaseURL, "/"), client: client}, nil
}

func (u *Uploader) Upload(ctx context.Context, objectKey string, body io.Reader, contentType string) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      awsv2.String(u.bucket),
		Key:         awsv2.String(objectKey),
		Body:        body,
		ContentType: awsv2.String(contentType),
	})
	if err != nil {
		return "", err
	}
	if u.publicBaseURL == "" {
		return "", nil
	}
	return fmt.Sprintf("%s/%s", u.publicBaseURL, strings.TrimLeft(objectKey, "/")), nil
}
