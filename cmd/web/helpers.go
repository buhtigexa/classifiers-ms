package main

import (
	"net/http"
	"runtime/debug"
)

// serverError handles any internal server errors
// Che, if something explodes internally, this is where we handle that quilombo
func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// First log the error with full stack trace, re importante for debugging viste
	a.logger.Error(err.Error(), 
		"method", r.Method, 
		"url", r.URL.Path,
		"trace", string(debug.Stack()),
	)
	
	// Then tell the user something went wrong, but not too much detail eh
	a.errorResponse(w, r, http.StatusInternalServerError, "sorry che, we had a problem internally")
}

// notFoundError handles 404 not found responses
// This is for when someone looks for something that no existe, viste?
func (a *application) notFoundError(w http.ResponseWriter, r *http.Request, id string) {
	a.logger.Error("Resource not found",
		"method", r.Method,
		"url", r.URL.Path,
		"id", id,
	)
	a.errorResponse(w, r, http.StatusNotFound, "che, we couldn't find what you're looking for")
}
