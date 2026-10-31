// Package db owns PostgreSQL infrastructure. Domain packages must not import it.
package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration")
	}
	cfg.MaxConns = 10
	cfg.MinConns = 0
	// No SQL/driver tracing: URLs and query arguments must not enter application logs.
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("database pool initialization failed")
	}
	return pool, nil
}
