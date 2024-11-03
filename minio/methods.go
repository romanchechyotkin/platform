package minio

import (
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"
)

type SaveObjectOptions struct {
	FileName    string
	FilePath    string
	BucketName  string
	ContentType string
}

type RemoveObjectOptions struct {
	FileName   string
	BucketName string
}

func (c *Client) SaveObject(ctx context.Context, opts *SaveObjectOptions) error {
	ctx, span := c.tracer.Start(ctx, "save-object")
	defer span.End()

	info, err := c.client.FPutObject(ctx, opts.BucketName, opts.FileName, opts.FilePath, minio.PutObjectOptions{
		ContentType: opts.ContentType,
	})
	if err != nil {
		c.log.Error("failed to upload object", err, slog.String("filename", opts.FileName))
		return err
	}

	c.log.Info("successfully uploaded object of size", slog.String("filename", opts.FileName), slog.Int64("size", info.Size))

	return nil
}

func (c *Client) RemoveObject(ctx context.Context, opts *RemoveObjectOptions) error {
	ctx, span := c.tracer.Start(ctx, "remove-object")
	defer span.End()

	err := c.client.RemoveObject(ctx, opts.BucketName, opts.FileName, minio.RemoveObjectOptions{})
	if err != nil {
		c.log.Error("failed to remove object", err, slog.String("filename", opts.FileName))
		return err
	}

	c.log.Info("successfully removed object", slog.String("filename", opts.FileName))

	return nil
}
