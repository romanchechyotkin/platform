package minio

import (
	"fmt"

	"github.com/TakeAway-Inc/platform/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	log    *logger.Logger
	tracer trace.Tracer
	client *minio.Client
}

// NewClient creates new minio client
func NewClient(log *logger.Logger, cfg *Config) (*Client, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	accessKeyID := cfg.User
	secretAccessKey := cfg.Password

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		log:    log,
		client: minioClient,
		tracer: otel.GetTracerProvider().Tracer("minio"),
	}, nil
}
