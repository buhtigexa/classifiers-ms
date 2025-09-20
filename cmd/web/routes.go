package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	// Aca definimos todas las routes, super important stuff
	mux := http.NewServeMux()
	
	// Home endpoint, nothing fancy viste
	mux.HandleFunc("GET /", app.Home)
	
	// CRUD operations for our classifiers, re piola
	mux.HandleFunc("POST /classifiers/create", app.CreateClassifier)
	mux.HandleFunc("GET /classifiers", app.ListClassifiers)
	mux.HandleFunc("GET /classifiers/{id}", app.GetClassifier)
	
	// Metrics endpoint for cuando everything explota
	mux.HandleFunc("GET /debug/metrics", app.metricsHandler)

	// Add the gzip middleware porque performance viste
	// This makes everything mas rapido, trust me
	handler := app.gzipMiddleware(mux)
	return handler
}
