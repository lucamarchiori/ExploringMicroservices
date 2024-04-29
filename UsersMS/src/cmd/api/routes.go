package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := mux.NewRouter()

	// Register custom NotFound and MethodNotAllowed handlers for the router.
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/user", app.createUserHandler).Methods("POST")
	router.HandleFunc("/user/{id}", app.showUserHandler).Methods("GET")
	router.HandleFunc("/user", app.indexUserHandler).Methods("GET")
	//router.HandleFunc("/users/{id}", updatePost).Methods("PUT")
	//router.HandleFunc("/users/{id}", deletePost).Methods("DELETE")
	return app.enableCORS(router)
}
