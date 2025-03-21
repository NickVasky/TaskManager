package mw

import (
	"TaskManager/repo"
	"TaskManager/service"
	"context"
	"log"
	"net/http"
)

type ctxKey string

const SessionKey ctxKey = "session"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func Recovery(f http.HandlerFunc) http.HandlerFunc {
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

func Auth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var session *repo.SessionEntity
		if s, err := service.SessionAuth(r); err != nil {
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
