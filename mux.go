// mux.go
package main

import (
	"net/http"
)

type Mux interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
	Handle(string, http.Handler)
	Handler(*http.Request) (http.Handler, string)
	ServeHTTP(http.ResponseWriter, *http.Request)
}
