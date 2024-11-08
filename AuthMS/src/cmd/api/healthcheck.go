package main

import (
	"net/http"
)

// Handle the healthcheck request
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which holds the information that we want to send in the response.
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"healthcheck": data}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
