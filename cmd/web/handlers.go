package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	t, err := app.tournaments.Latest()
	if err != nil {
		w.Write([]byte("Something went wrong...")) // TODO
		return
	}

	for _, tournament := range t {
		fmt.Fprintf(w, "%v\n", tournament)
	}
}
