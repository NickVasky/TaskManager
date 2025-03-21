package common

import "net/http"

type Route struct {
	Pattern string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}
