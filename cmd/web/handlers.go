package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/forms"
	"github.com/julienschmidt/httprouter"
)

const layout = "2006-01-02T15:04"

func (app *application) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := app.tournaments.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.gohtml", &templateData{
		Tournaments: t,
	})
}

func (app *application) createTournamentForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := url.Values{}
	v.Set("datetime", time.Now().Format(layout))
	form := forms.New(v)

	app.render(w, r, "create.page.gohtml", &templateData{
		Form: form,
	})
}

func (app *application) createTournament(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "datetime", "type")
	form.MaxLength("title", 100)
	form.PermittedValues("type", "0", "1", "2", "3")
	form.ValidDate("datetime", layout)

	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{
			Form: form,
		})
		return
	}

	_, err = app.tournaments.Insert(form.Get("title"), form.Get("datetime"), form.Get("type"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// TODO: redirect to the newly created tournament
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
