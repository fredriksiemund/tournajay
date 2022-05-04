package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/db"
	"github.com/fredriksiemund/tournament-planner/pkg/models/mongodb"
	"github.com/fredriksiemund/tournament-planner/pkg/models/postgres"
	"github.com/golangcollege/sessions"
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
	users           *mongodb.UserModel
}

func main() {
	// Parsing the runtime configuration settings for the application
	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "postgres://root:root@localhost:5432/tournajay", "PostgreSQL connection string")
	connStrMongo := flag.String("connStrMongo", "mongodb://root:root@localhost:27017", "MongoDb connection string")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session secret key")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Establish DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	psql, err := db.PostgresConnect(*connStr, ctx)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer psql.Close(ctx)

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.MongoConnect(*connStrMongo, ctx)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer client.Disconnect(ctx)
	mdb := client.Database("tournajay")

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
		games:           &postgres.GameModel{Db: psql},
		infoLog:         infoLog,
		participants:    &postgres.ParticipantModel{Db: psql},
		session:         session,
		teams:           &postgres.TeamModel{Db: psql},
		templateCache:   templateCache,
		tournaments:     &postgres.TournamentModel{Db: psql},
		tournamentTypes: &postgres.TournamentTypeModel{Db: psql},
		users:           &mongodb.UserModel{Coll: mdb.Collection("users")},
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
