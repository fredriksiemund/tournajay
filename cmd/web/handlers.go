package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/forms"
	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/idtoken"
)

const layout = "2006-01-02T15:04"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	t, err := app.tournaments.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.gohtml", &templateData{
		Tournaments: t,
	})
}

func (app *application) createTournamentForm(w http.ResponseWriter, r *http.Request) {
	v := url.Values{}
	v.Set("datetime", time.Now().Format(layout))
	form := forms.New(v)

	app.render(w, r, "create.page.gohtml", &templateData{
		Form: form,
	})
}

func (app *application) createTournament(w http.ResponseWriter, r *http.Request) {
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

	id, err := app.tournaments.Insert(form.Get("title"), form.Get("description"), form.Get("datetime"), form.Get("type"), app.getCurrentUser(r).Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("tournament/%d", id), http.StatusSeeOther)
}

func (app *application) showTournament(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	t, err := app.tournaments.One(id)
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

func (app *application) showGamePlan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	t, err := app.tournaments.One(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "game-plan.page.gohtml", &templateData{
		Tournament: t,
	})
}

func (app *application) removeTournament(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	token := r.PostForm.Get("credential")
	payload, err := idtoken.Validate(context.Background(), token, googleClientId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id := fmt.Sprintf("%v", payload.Claims["sub"])
	name := fmt.Sprintf("%v", payload.Claims["name"])
	email := fmt.Sprintf("%v", payload.Claims["email"])
	picture := fmt.Sprintf("%v", payload.Claims["picture"])

	err = app.users.Insert(id, name, email, picture)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, sessionKeyIdToken, token)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, sessionKeyIdToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
