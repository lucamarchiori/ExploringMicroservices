package main

import (
	"errors"
	"log"
	"lucamarchiori/MicroserviceBoilerplate/internal/data"
	"lucamarchiori/MicroserviceBoilerplate/internal/validator"

	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	// Anonymous struct to hold the information expected to be in the HTTP request body
	var input struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"`
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

	user.Password, err = data.HashPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Validation
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Model insert
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Printf("Index users handler called")

	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.Get(id)

	if err != nil {
		app.logger.Printf("Error: %v", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) indexUsersHandler(w http.ResponseWriter, r *http.Request) {

	// Request inputs (filters)
	var input struct {
		Email string `json:"email"`
		data.Filters
	}

	// Read parameters from request
	qs := r.URL.Query()
	input.Email = app.readString(qs, "email", "")

	// Store the query result
	var users []*data.User
	users, err := app.models.Users.Index(input.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return response
	err = app.writeJSON(w, http.StatusOK, envelope{"users": users, "input": input}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
