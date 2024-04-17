package data

import (
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
