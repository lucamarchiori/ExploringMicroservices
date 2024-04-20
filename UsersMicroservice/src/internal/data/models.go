package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. Return this from Get() method when looking up a model that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Models struct which wraps the Model.
type Models struct {
	Users UserModel
}

// Method which returns a Models struct containing the initialized Model.
func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
