package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func PostgresConnect(connStr string, ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
