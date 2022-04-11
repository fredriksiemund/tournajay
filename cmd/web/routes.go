package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()

	// Middleware
	r.Use(app.recoverPanic)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	// Public routes
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Authenticated Routes
	r.Group(func(r chi.Router) {
		r.Use(app.session.Enable)
		r.Use(app.authenticate)

		r.Get("/", app.home)
		r.Route("/tournament", func(r chi.Router) {
			r.Post("/", app.createTournament)
			r.Get("/", app.createTournamentForm)
			r.Get("/{id}", app.showTournament)
			r.Delete("/{id}", app.removeTournament)
		})
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", app.loginUser)
			r.Post("/logout", app.logoutUser)
		})
	})

	return r
}
