package data

import (
	"database/sql"
	"lucamarchiori/MicroserviceBoilerplate/internal/validator"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name,omitempty"`
	Surname   string    `json:"surname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(user.Surname != "", "surname", "must be provided")
	v.Check(user.Email != "", "email", "must be provided")
}

// Define a struct type which wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the table.
func (m UserModel) Insert(user *User) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the table.
func (m UserModel) Get(id int64) (*User, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the table.
func (m UserModel) Update(user *User) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the table.
func (m UserModel) Delete(id int64) error {
	return nil
}
