package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := mux.NewRouter()

	// Middleware to log incoming request URI
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the request URI
			app.logger.Println(r.RequestURI)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	})

	// Register custom NotFound and MethodNotAllowed handlers for the router.
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// API Routes and handlers
	router.HandleFunc("/users-ms/healthcheck", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/users-ms/users", app.createUserHandler).Methods("POST")
	router.HandleFunc("/users-ms/users/{id}", app.showUserHandler).Methods("GET")
	router.HandleFunc("/users-ms/users", app.indexUsersHandler).Methods("GET")
	//router.HandleFunc("/users/{id}", updatePost).Methods("PUT")
	//router.HandleFunc("/users/{id}", deletePost).Methods("DELETE")
	return app.enableCORS(router)
}
