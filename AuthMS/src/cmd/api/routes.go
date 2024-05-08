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

	// Custom 404 and 405 handlers
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// Endpoints
	router.HandleFunc("/auth-ms/healthcheck", app.healthcheckHandler).Methods("GET")
	router.HandleFunc("/auth-ms/auth/login", app.loginHandler).Methods("POST")
	router.HandleFunc("/auth-ms/auth/signup", app.signUpHandler).Methods("POST")
	return app.enableCORS(router)
}
