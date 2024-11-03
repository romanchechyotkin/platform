package postgresql

import (
	"context"
	"log/slog"
	"time"

	"github.com/TakeAway-Inc/platform/logger"

	"go.uber.org/fx"
)

const moduleName = "postgres"

func NewModule() fx.Option {
	return fx.Module(
		moduleName,

		fx.Provide(
			func(log *logger.Logger, cfg *Config) (*Postgres, error) {
				return New(log, cfg, ConnAttempts(20),
					MinPoolSize(5),
					MaxConnLifetime(time.Hour),
					MaxConnIdleTime(30*time.Minute),
					ConnHealthCheckPeriod(time.Minute),
				)
			},
		),

		fx.Invoke(func(
			lc fx.Lifecycle,
			p *Postgres,
		) {
			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						return nil
					},
					OnStop: func(_ context.Context) error {
						p.Close()
						return nil
					},
				},
			)
		}),

		fx.Decorate(func(log *logger.Logger) *logger.Logger {
			return log.With(slog.String("module", moduleName))
		}),
	)
}
