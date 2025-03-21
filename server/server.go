package server

import (
	"TaskManager/server/api"
	"TaskManager/server/pages"
	"fmt"
	"net/http"
)

func BuildMux() *http.ServeMux {
	routes := append(api.Routes, pages.Routes...)
	mux := http.NewServeMux()
	for _, r := range routes {
		pattern := fmt.Sprintf("%s %s", r.Method, r.Pattern)
		mux.HandleFunc(pattern, r.Handler)
	}
	return mux
}

func Serve() {
	mux := BuildMux()
	fs := http.FileServer(http.Dir("templates/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", mux)
}
