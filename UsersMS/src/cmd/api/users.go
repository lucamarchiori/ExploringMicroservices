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

	err = user.Password.Set(input.Password)
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
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := data.User{
		Id:      int(id),
		Name:    "Luca",
		Surname: "Marchiori",
		Email:   "ladmin@example.com",
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) indexUsersHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Printf("Index users handler called")
	var users []*data.User

	users, err := app.models.Users.Index()
	if err != nil {
		app.logger.Printf("Error: %v", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
}
