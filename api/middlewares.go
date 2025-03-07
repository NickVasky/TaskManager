package api

import (
	"TaskManager/repo"
	"context"
	"log"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type contextKey string

const SessionKey contextKey = "session"

func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func mw_auth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var session *repo.SessionEntity
		if s, err := Authorize(r); err != nil {
			err := http.StatusUnauthorized
			http.Error(w, "Unauthorized", err)
			return
		} else {
			session = s
		}
		ctx := context.WithValue(r.Context(), SessionKey, session)
		f(w, r.WithContext(ctx))
	}
}

func mw_methodCheck(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != patternMethods[r.Pattern] {
			err := http.StatusMethodNotAllowed
			http.Error(w, "Invalid Method", err)
			return
		}
		f(w, r)
	}
}

func mw_errorRecovery(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := "Caught panic: %v, Stack trace: %s"
				log.Printf(msg, err)

				err := http.StatusInternalServerError
				http.Error(w, "Internal Server Error", err)
			}
		}()
		f(w, r)
	}
}
