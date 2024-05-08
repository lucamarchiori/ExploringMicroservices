package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lucamarchiori/MicroserviceBoilerplate/internal/data"
	"net/http"
	"time"
)

// Handle the login request
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get the email and hashed password from the request body
	// Step 2: get the corresponding user from the users microservice
	// Step 3: if the user is not found, return a 401 Unauthorized response
	// Step 4: if the user is found, egenerate a new bearer token
	// Step 5: save the new token in the tokens table
	// Step 6: return the token to the client

	// Get the email and password from the request body
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Error reading JSON")
		app.badRequestResponse(w, r, err)
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get the corresponding user from the users microservice
	app.logger.Println("Sending request to users microservice")
	httpClient := &http.Client{}
	endpoint := "http://" + app.services["users"] + ":4000" + "/users?email=" + input.Email
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Println(req)
	res, err := httpClient.Do(req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer res.Body.Close()

	app.logger.Println(res)

	// Step 3: if the user is not found, return a 401 Unauthorized response
	if res.StatusCode == http.StatusNotFound {
		app.unauthorizedResponse(w, r, errors.New("the users microservice could not find the user with the specified email address"))
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New("received a non-200 status code from the users microservice")
		app.serverErrorResponse(w, r, err)
		return
	}

	// Read response body
	body, err := io.ReadAll(res.Body)

	// Step 4: if the user is found, save it in the user variable
	type Users []struct {
		Id        int       `json:"id"`
		Name      string    `json:"name,omitempty"`
		Surname   string    `json:"surname,omitempty"`
		Email     string    `json:"email,omitempty"`
		Password  string    `json:"password,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	var response struct {
		Users Users `json:"users,omitempty"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	if len(response.Users) <= 0 {
		app.logger.Println("No users found with the specified email address")
		app.unauthorizedResponse(w, r, errors.New("the users microservice could not find the user with the specified email address"))
	}

	// Get the first user returned and check for password match
	pwdMatch, err := data.PasswordMatch(input.Password, response.Users[0].Password)

	if !pwdMatch {
		app.logger.Println("Password does not match")
		app.unauthorizedResponse(w, r, errors.New("wrong password"))
	}

	// Generate a new bearer token
	app.logger.Println("Password matches, generating new token")
	token := "Bearer " + data.GenerateSecureToken(32)

	// Return the token to the client along with the user
	err = app.writeJSON(w, http.StatusOK, envelope{"user": response.Users[0], "token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// Handle the signup request
func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {

}
