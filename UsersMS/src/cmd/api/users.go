package main

import (
	"fmt"
	"log"
	"lucamarchiori/MicroserviceBoilerplate/internal/data"
	"lucamarchiori/MicroserviceBoilerplate/internal/validator"

	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	// Anonymous struct to hold the information expected to be in the HTTP request body
	var input struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email"`
	}

	// json.Decoder instance which reads from the request body, and then use the Decode() method to decode the body contents into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(r.Body)
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:    input.Name,
		Surname: input.Surname,
		Email:   input.Email,
	}

	// Validation
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := data.User{
		ID:       int(id),
		Name:     "Luca",
		Surname:  "Marchiori",
		Email:    "ladmin@example.com",
		Password: "password",
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) indexUserHandler(w http.ResponseWriter, r *http.Request) {
	users := []*data.User{
		{
			ID:       1,
			Name:     "Luca",
			Surname:  "Marchiori",
			Email:    "luca.marchiori@example.com",
			Password: "password",
		}}

	err := app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
