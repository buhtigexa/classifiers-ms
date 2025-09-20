package main

import (
	"net/http"
)

func (app *application) metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := app.metrics.GetMetrics()
	err := app.writeJSON(w, http.StatusOK, envelope{
		"metrics": map[string]interface{}{
			"open_connections":     metrics.OpenConnections,
			"in_use_connections":   metrics.InUseConnections,
			"wait_count":          metrics.WaitCount,
			"max_idle_closed":     metrics.MaxIdleTimeClosed,
		},
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}