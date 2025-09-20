package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (a *application) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Classifier Home Page!"))
}

func (a *application) CreateClassifier(w http.ResponseWriter, r *http.Request) {
	a.model.Insert("Example Classifier This is an example classifier.")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Classifier created successfully!"))
	w.Header().Add("Server", "ClassifierServer/1.0")

}

func (a *application) GetClassifier(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.serverError(w, r, fmt.Errorf("method not allowed"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Server", "ClassifierServer/1.0")
	w.Write([]byte("Here is your classifier!"))
}

func (a *application) ListClassifiers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := strconv.Atoi(id) // Just to avoid unused import error
	if err != nil || i < 1 {
		a.notFoundError(w, r, id)
		return
	}
	w.Header().Add("Server", "ClassifierServer/1.0")
	w.Write([]byte(fmt.Sprintf("Listing classifier with ID: %s", id)))
}
