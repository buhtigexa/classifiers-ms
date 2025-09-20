//go:build docs
// +build docs

// Package main Classifier API
//
// Che boludo, this is our super optimized classifier API documentation
// We've got some zarpado endpoints here for managing classifiers
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
// Contact: Classifier Support<support@buhtigexa.net>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import "classifier.buhtigexa.net/internal/models"

// swagger:route POST /classifiers/create classifiers createClassifier
// Create a new classifier, super simple viste?
// responses:
//   201: classifierResponse
//   400: errorResponse

// swagger:parameters createClassifier
type createClassifierParams struct {
	// in: body
	// required: true
	Body struct {
		// The name of the classifier
		// required: true
		Name string `json:"name"`
		// An optional description
		Description string `json:"description,omitempty"`
		// Whether the classifier is active
		IsActive *bool `json:"is_active,omitempty"`
	}
}

// swagger:response classifierResponse
type swaggerClassifierResponse struct {
	// in: body
	Body struct {
		Classifier *models.Classifier `json:"classifier"`
	}
}

// swagger:response errorResponse
type swaggerErrorResponse struct {
	// in: body
	Body struct {
		Error string `json:"error"`
	}
}

// swagger:route GET /classifiers/{id} classifiers getClassifier
// Get a classifier by ID
// responses:
//   200: classifierResponse
//   404: errorResponse

// swagger:parameters getClassifier
type getClassifierParams struct {
	// The ID of the classifier
	// in: path
	// required: true
	ID int64 `json:"id"`
}

// swagger:route GET /classifiers classifiers listClassifiers
// List all classifiers with pagination
// responses:
//   200: listResponse
//   400: errorResponse

// swagger:parameters listClassifiers
type listClassifiersParams struct {
	// The page number
	// in: query
	// minimum: 1
	// default: 1
	Page int `json:"page"`
	
	// Items per page
	// in: query
	// minimum: 1
	// maximum: 100
	// default: 20
	PageSize int `json:"page_size"`
}

// swagger:response listResponse
type swaggerListResponse struct {
	// in: body
	Body struct {
		Classifiers []*models.Classifier `json:"classifiers"`
		Metadata    struct {
			Total    int `json:"total"`
			Page     int `json:"page"`
			PageSize int `json:"page_size"`
			Pages    int `json:"pages"`
		} `json:"metadata"`
	}
}