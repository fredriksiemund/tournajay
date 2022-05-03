package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/models/postgres"
	"github.com/golangcollege/sessions"
	"github.com/jackc/pgx/v4"
)

type contextKey string

const contextKeyCurrentUser = contextKey("isAuthenticated")
const sessionKeyIdToken = "idToken"
const googleClientId = "879593153148-6pho9arld8k17qol30c23hlr02i8qeru.apps.googleusercontent.com"

type application struct {
	errorLog        *log.Logger
	games           *postgres.GameModel
	infoLog         *log.Logger
	participants    *postgres.ParticipantModel
	session         *sessions.Session
	teams           *postgres.TeamModel
	templateCache   map[string]*template.Template
	tournaments     *postgres.TournamentModel
	tournamentTypes *postgres.TournamentTypeModel
	users           *postgres.UserModel
}

func main() {
	// Parsing the runtime configuration settings for the application
	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "postgres://root:root@localhost:5432/tournajay", "PostgreSQL connection string")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session secret key")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Establish DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgresConnect(*connStr, ctx)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close(ctx)

	// Set up template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Establishing the dependencies for the handlers (depenency injection)
	app := &application{
		errorLog:        errorLog,
		games:           &postgres.GameModel{Db: db},
		infoLog:         infoLog,
		participants:    &postgres.ParticipantModel{Db: db},
		session:         session,
		teams:           &postgres.TeamModel{Db: db},
		templateCache:   templateCache,
		tournaments:     &postgres.TournamentModel{Db: db},
		tournamentTypes: &postgres.TournamentTypeModel{Db: db},
		users:           &postgres.UserModel{Db: db},
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

func postgresConnect(connStr string, ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	// Since connections to the database are established lazily,
	// we can verify that everything is set up correctly by calling db.Ping()
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
