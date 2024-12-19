package api

import "net/http"

type API interface {
	Start()
	// ServeHTTP is used in component test
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
