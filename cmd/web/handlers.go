package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := app.tournaments.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.gohtml", &templateData{Tournaments: t})
}

func (app *application) createTournamentForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := app.tournamentTypes.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "create.page.gohtml", &templateData{TournamentTypes: t})
}
