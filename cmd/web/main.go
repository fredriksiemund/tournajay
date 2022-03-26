package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/fredriksiemund/tournament-planner/pkg/models/postgres"
	"github.com/jackc/pgx/v4"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	tournaments *postgres.TournamentModel
}

func main() {
	// Parsing the runtime configuration settings for the application
	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "postgres://root:root@localhost:5432/dev", "PostgreSQL connection string")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Establish DB connection
	db, err := openDb(*connStr)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close(context.Background())

	// Establishing the dependencies for the handlers (depenency injection)
	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		tournaments: &postgres.TournamentModel{Db: db},
	}

	// Running the HTTP server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDb(connStr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	// Since connections to the database are established lazily,
	// we can verify that everything is set up correctly by calling db.Ping()
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}
