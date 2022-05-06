package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/db"
)

func main() {
	down := flag.Bool("down", false, "Clear database")
	connStr := flag.String("connStr", "postgres://root:root@localhost:5432/tournajay", "PostgreSQL connection string")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	psql, err := db.PostgresConnect(*connStr, ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer psql.Close(ctx)

	if *down {
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_, err := psql.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		path := filepath.Join("sql", "data.sql")
		c, ioErr := ioutil.ReadFile(path)
		if ioErr != nil {
			log.Fatal(err)
		}
		sql := string(c)

		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_, err := psql.Exec(ctx, sql)
		if err != nil {
			log.Fatal(err)
		}
	}
}
