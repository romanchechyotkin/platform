package postgresql

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func MinPoolSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func ConnAttempts(size int) Option {
	return func(p *Postgres) {
		p.connAttempts = size
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

func MaxConnLifetime(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.maxConnLifetime = timeout
	}
}

func MaxConnIdleTime(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.maxConnIdleTime = timeout
	}
}

func ConnHealthCheckPeriod(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.healthCheckPeriod = timeout
	}
}
