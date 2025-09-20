package main

import (
	"encoding/json"
	"net/http"

	"classifier.buhtigexa.net/internal/models"
)

// Che, these structs are pre-defined to avoid that reflection thing during runtime
// I mean, its faster this way, viste?
type envelope map[string]interface{}

type listResponse struct {
	Classifiers []*models.Classifier `json:"classifiers"`
	Metadata    listMetadata        `json:"metadata"`
}

type listMetadata struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Pages    int `json:"pages"`
}

type classifierResponse struct {
	Classifier *models.Classifier `json:"classifier"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logger.Error("Error writing response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// badRequestError returns a 400 Bad Request response with the error message
// Che, this one is for when the user sends us cualquier cosa
func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
