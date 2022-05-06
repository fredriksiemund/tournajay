package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/fredriksiemund/tournament-planner/pkg/forms"
	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/fredriksiemund/tournament-planner/pkg/tournaments"
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
	v.Set("date", time.Now().Format(layout))
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
	form.Required("title", "date", "type")
	form.MaxLength("title", 100)
	form.MaxLength("description", 1000)
	form.PermittedValues("type", "1", "2", "3", "4")
	form.ValidDate("date", layout)

	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.tournaments.Insert(form.Get("title"), form.Get("description"), form.Get("date"), form.Get("type"), app.getCurrentUser(r).Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/tournament/%d", id), http.StatusSeeOther)
}

func (app *application) showTournament(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tournamentId"))
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

	app.render(w, r, "general.page.gohtml", &templateData{
		Tournament: t,
	})
}

func (app *application) showSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tournamentId"))
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

	teams, err := app.teams.All(t.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	games, err := app.games.All(t.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	rounds := tournaments.SingleEliminationTemplate(games, teams)

	app.render(w, r, "schedule.page.gohtml", &templateData{
		Teams:      teams,
		Tournament: t,
		Rounds:     rounds,
	})
}

func (app *application) createSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tournamentId"))
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

	// if t.Creator.Id != app.getCurrentUser(r).Id {
	// 	app.clientError(w, http.StatusForbidden)
	// 	return
	// }

	// Check if a game already exists
	exists, err := app.games.Exists(id)
	if err != nil {
		app.serverError(w, err)
		return
	} else if exists {
		app.clientError(w, http.StatusConflict)
		return
	}

	// Generate teams
	nbrOfTeams := int(math.Ceil(float64(len(t.Participants)) / 2.0))
	teamIds, err := app.teams.Insert(nbrOfTeams, t.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Asign teams
	err = app.participants.AssignTeams(t.Id, t.Participants, teamIds)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Generate schedule
	bracket, err := tournaments.NewSingleElimination(teamIds)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Insert schedule into database
	err = app.games.InsertSingleEliminationGames(t.Id, bracket)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/tournament/%d/schedule", id), http.StatusSeeOther)
}

func (app *application) removeTournament(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tournamentId"))
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

	err = app.users.Upsert(id, name, email, picture)
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

func (app *application) createParticipant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tournamentId"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.participants.Insert(id, app.getCurrentUser(r).Id)
	if err != nil && !errors.Is(err, models.ErrDuplicate) {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/tournament/%d", id), http.StatusSeeOther)
}

func (app *application) removeParticipant(w http.ResponseWriter, r *http.Request) {
	tournamentId, err := strconv.Atoi(chi.URLParam(r, "tournamentId"))
	if err != nil || tournamentId < 1 {
		app.notFound(w)
		return
	}

	userId := chi.URLParam(r, "userId")
	if utf8.RuneCountInString(userId) < 1 {
		app.notFound(w)
		return
	}

	err = app.participants.Delete(tournamentId, userId)
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
