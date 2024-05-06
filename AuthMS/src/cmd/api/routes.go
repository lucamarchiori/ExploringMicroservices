package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := mux.NewRouter()

	// Custom 404 and 405 handlers
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// Endpoints
	router.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/auth/login", app.loginHandler).Methods("POST")
	router.HandleFunc("/auth/signup", app.signUpHandler).Methods("POST")
	return app.enableCORS(router)
}
