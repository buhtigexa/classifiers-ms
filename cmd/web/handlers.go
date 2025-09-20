package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"classifier.buhtigexa.net/internal/models"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{
		"message": "Welcome to the Classifier API",
		"status":  "available",
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

type createClassifierRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

func (app *application) CreateClassifier(w http.ResponseWriter, r *http.Request) {
	var req createClassifierRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if req.Name == "" {
		app.badRequestError(w, r, fmt.Errorf("name is required"))
		return
	}

	// Get the description string value or empty string if nil
	var description string
	if req.Description != nil {
		description = *req.Description
	}

	id, err := app.model.Insert(req.Name, description, req.IsActive)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create response with all fields
	err = app.writeJSON(w, http.StatusCreated, envelope{
		"classifier": map[string]interface{}{
			"id":          id,
			"name":        req.Name,
			"description": req.Description,
			"is_active":   req.IsActive,
		},
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) GetClassifier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestError(w, r, fmt.Errorf("invalid id parameter"))
		return
	}

	classifier, err := app.model.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundError(w, r, fmt.Sprintf("%d", id))
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"classifier": classifier}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) ListClassifiers(w http.ResponseWriter, r *http.Request) {
	// Bueno, aca parseamos los params de paginacion
	// Es importante porque si no limitamos esto, se va todo al carajo
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		var err error
		page, err = strconv.Atoi(p)
		if err != nil || page < 1 {
			// Mandaron cualquier fruta en el page parameter
			app.badRequestError(w, r, fmt.Errorf("invalid page parameter"))
			return
		}
	}

	pageSize := 20
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		var err error
		pageSize, err = strconv.Atoi(ps)
		if err != nil || pageSize < 1 || pageSize > 100 {
			app.badRequestError(w, r, fmt.Errorf("invalid page_size parameter"))
			return
		}
	}

	classifiers, total, err := app.model.List(models.ListClassifiersOptions{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	response := listResponse{
		Classifiers: classifiers,
		Metadata: listMetadata{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Pages:    (total + pageSize - 1) / pageSize,
		},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": response}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
