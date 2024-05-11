package data

import (
	"errors"
)

// Define a custom ErrRecordNotFound error. Return this from Get() method when looking up a model that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Models struct which wraps the Model.
type Models struct {
	Auth AuthModel
}

// Method which returns a Models struct containing the initialized Model.
func NewModels() Models {
	return Models{
		Auth: AuthModel{},
	}
}
