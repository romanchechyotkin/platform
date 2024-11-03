package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/TakeAway-Inc/platform/logger"

	"github.com/Masterminds/squirrel"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxPoolSize       = 1
	defaultMinPoolSize       = 1
	defaultConnAttempts      = 10
	defaultConnTimeout       = time.Second
	defaultConnLifetime      = time.Minute
	defaultIdleTime          = time.Minute
	defaultHealthCheckPeriod = time.Minute
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Stat() *pgxpool.Stat
	Ping(ctx context.Context) error
}

type Postgres struct {
	maxPoolSize       int
	minPoolSize       int
	connAttempts      int
	maxConnLifetime   time.Duration
	maxConnIdleTime   time.Duration
	healthCheckPeriod time.Duration
	connTimeout       time.Duration

	Builder squirrel.StatementBuilderType
	Pool    PgxPool
	Log     *logger.Logger
}

func New(log *logger.Logger, cfg *Config, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:       defaultMaxPoolSize,
		minPoolSize:       defaultMinPoolSize,
		connAttempts:      defaultConnAttempts,
		maxConnLifetime:   defaultConnLifetime,
		maxConnIdleTime:   defaultIdleTime,
		healthCheckPeriod: defaultHealthCheckPeriod,
		connTimeout:       defaultConnTimeout,
		Builder:           squirrel.StatementBuilderType{},
		Log:               log,
		Pool:              nil,
	}

	for _, opt := range opts {
		opt(pg)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	log.Debug("connection url", slog.String("url", url))

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	pgConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgConfig.MaxConns = int32(pg.maxPoolSize)         // Maximum number of connections in the pool
	pgConfig.MinConns = int32(pg.minPoolSize)         // Minimum number of connections in the pool
	pgConfig.MaxConnLifetime = pg.maxConnLifetime     // Maximum connection lifetime
	pgConfig.MaxConnIdleTime = pg.maxConnIdleTime     // Maximum idle time before connection is closed
	pgConfig.HealthCheckPeriod = pg.healthCheckPeriod // Health check period
	pgConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), pgConfig)
		if err == nil {
			break
		}

		log.Debug("postgres is trying to connect", slog.Int("attempts left", pg.connAttempts))
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if cfg.AutoCreate {
		exists, err := CreateDatabase(ctx, url)
		if err != nil {
			log.Error("failed to create database", err)

			return nil, err
		}

		if exists {
			log.Info("the database already exists", slog.String("database", cfg.Database))
		} else {
			log.Info("the database was created successfully", slog.String("database", cfg.Database))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("pgdb - New - pgxpool.ConnectConfig: %w", err)
	}

	return pg, pg.Pool.Ping(ctx)
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
