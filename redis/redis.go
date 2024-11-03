package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/TakeAway-Inc/platform/logger"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Database int
}

type Client struct {
	log *logger.Logger

	redisClient *redis.Client
}

func NewClient(log *logger.Logger, cfg *Config) (*Client, error) {
	log = log.With(slog.String("module", "redis"))

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DB:   cfg.Database,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return nil, err
	}

	return &Client{
		log:         log,
		redisClient: redisClient,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.log.Info("set value into redis", slog.String("key", key), slog.Any("value", value))
	return c.redisClient.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.redisClient.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.redisClient.Del(ctx, key).Err()
}

func (c *Client) Close() error {
	return c.redisClient.Close()
}
