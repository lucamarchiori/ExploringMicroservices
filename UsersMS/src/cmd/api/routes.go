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

	// API Routes and handlers
	router.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/users", app.createUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", app.showUserHandler).Methods("GET")
	router.HandleFunc("/users", app.indexUsersHandler).Methods("GET")
	//router.HandleFunc("/users/{id}", updatePost).Methods("PUT")
	//router.HandleFunc("/users/{id}", deletePost).Methods("DELETE")
	return app.enableCORS(router)
}
