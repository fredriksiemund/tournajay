package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) getCurrentUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyCurrentUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentUser = app.getCurrentUser(r)
	td.RequestURI = r.URL.RequestURI()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.gohtml'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(w, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}
}
