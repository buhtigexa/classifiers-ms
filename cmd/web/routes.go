package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.Home)
	mux.HandleFunc("POST /create", app.CreateClassifier)
	mux.HandleFunc("/list", app.ListClassifiers)
	mux.HandleFunc("/classifier/{id}", app.GetClassifier)
	return mux

}
