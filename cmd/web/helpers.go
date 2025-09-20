package main

import (
	"net/http"
	"runtime/debug"
)

func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error(err.Error(), "method", r.Method, "url", r.URL.Path, debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (a *application) notFoundError(w http.ResponseWriter, r *http.Request, id string) {
	a.logger.Error("method", r.Method, "url", r.URL.Path, "id", id)
	http.Error(w, "Not found", http.StatusNotFound)
}
