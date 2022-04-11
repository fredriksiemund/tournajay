package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"google.golang.org/api/idtoken"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if id exists
		exists := app.session.Exists(r, "idToken")
		if !exists {
			// Proceed to next handler without setting context
			next.ServeHTTP(w, r)
			return
		}

		// Check if a user with the provided id exists
		token := app.session.GetString(r, "idToken")
		payload, err := idtoken.Validate(context.Background(), token, googleClientId)
		if err != nil {
			// Proceed to next handler without setting context
			next.ServeHTTP(w, r)
			return
		}

		id := fmt.Sprintf("%v", payload.Claims["sub"])
		_, err = app.users.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				// User does not exist, remove from session data and proceed without setting context
				app.session.Remove(r, "idToken")
				next.ServeHTTP(w, r)
				return
			} else {
				// Something else went wrong
				app.serverError(w, err)
				return
			}
		}

		// Include information in request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyIsAuthenticated, true)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
