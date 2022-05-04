package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoConnect(connStr string, ctx context.Context) (*mongo.Client, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return conn, nil
}

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
