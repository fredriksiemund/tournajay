package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/forms"
	"github.com/fredriksiemund/tournament-planner/pkg/models"
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
	form.MaxLength("description", 1000)
	form.PermittedValues("type", "1", "2", "3", "4")
	form.ValidDate("datetime", layout)

	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.tournaments.Insert(form.Get("title"), form.Get("description"), form.Get("datetime"), form.Get("type"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("tournament/%d", id), http.StatusSeeOther)
}

func (app *application) showTournament(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	t, err := app.tournaments.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.gohtml", &templateData{
		Tournament: t,
	})
}

func (app *application) removeTournament(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.tournaments.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
