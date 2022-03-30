package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.GET("/", app.home)
	router.GET("/tournament", app.createTournamentForm)
	router.POST("/tournament", app.createTournament)
	router.GET("/tournament/:id", app.showTournament)

	return router
}
