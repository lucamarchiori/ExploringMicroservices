package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Define an envelope type.

type envelope map[string]interface{}

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIdParam(r *http.Request) (int64, error) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// Retrieve the "email" URL parameter from the request
func (app *application) readEmailParam(r *http.Request) (string, error) {
	params := mux.Vars(r)
	email := params["email"]
	if email == "" {
		return "", errors.New("invalid email parameter")
	}
	return email, nil
}

// Helper for sending responses
// Takes http.ResponseWriter, the HTTP status code to send, the enveloped data to encode to JSON, and a
// header map containing any additional HTTP headers to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	// Encode the data to JSON, returning the error if there is one.
	js, err := json.MarshalIndent(data, "", "\t")
	// While using json.MarshalIndent() is positive from a readability and user-experience point of view,
	// it leads to overhead: takes 65% longer to run and uses around 30% more memory than json.Marshal()

	if err != nil {
		return err
	}
	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Add any headers that we want to include.
	// Loop through the header map and add each header to the http.ResponseWrite.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// Helper for reading JSON request bodies
// Takes http.ResponseWriter, the HTTP request to read the body from, and a pointer to the destination
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	// If the JSON from the client now includes any field which cannot be mapped to the target destination, the decoder will return an error instead of just ignoring the field.
	dec.DisallowUnknownFields()
	// Decode the request body into the target destination.
	err := dec.Decode(dst)

	// If there is an error during decoding, start the error triage
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error for syntax errors in the JSON.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		// Catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then include that in error message to make it easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
			// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
			// pointer to Decode(). We catch this and panic, rather than returning an error
			// to our handler. At the end of this chapter we'll talk about panicking
			// versus returning errors, and discuss why it's an appropriate thing to do in
			// this specific situation.

		// If the JSON contains a field which cannot be mapped to the target destination then
		// the decoder will return an error message in the format 'json: unknown field "<field name>"'.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds in size the decode will fail
		case errors.Is(err, &http.MaxBytesError{}):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// For anything else, return the error message as-is.
		default:
			return err
		}
	}
	return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	// Otherwise return the string.
	return s
}
