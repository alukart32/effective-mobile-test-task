// Package postgres provides pgxpool.Pool.
package postgres

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

// Get returns an instance of pgxpool.Pool.
func Get(dsn string) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {
		var cfg *pgxpool.Config

		cfg, err = conf(dsn)
		if err != nil {
			return
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
		if err != nil {
			return
		}

		// Ping a new pool.
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err = pool.Ping(ctx)
		if err != nil {
			return
		}
	})

	return pool, err
}

// conf prepares pgxpool.Config.
func conf(dsn string) (*pgxpool.Config, error) {
	opts := env.Options{RequiredIfNoDef: true}

	var cfg poolConf
	err := env.ParseWithOptions(&cfg, opts)
	if err != nil {
		return nil, err
	}

	if len(dsn) == 0 {
		return nil, errors.New("DSN is empty")
	}
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	conf.MaxConns = int32(cfg.MaxConns)

	return conf, nil
}

// poolConf is the representation of postgres pool settings.
type poolConf struct {
	MaxConns    int32         `env:"DATABASE_MAX_CONNS" envDefault:"5"`
	PingTimeout time.Duration `env:"PING_TIMEOUT" envDefault:"300ms"`
}
