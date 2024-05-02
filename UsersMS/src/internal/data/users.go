package data

import (
	"context"
	"database/sql"
	"errors"
	"lucamarchiori/MicroserviceBoilerplate/internal/validator"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name,omitempty"`
	Surname   string    `json:"surname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  password  `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type password struct {
	plaintext *string
	hash      []byte
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

// Calculates the bcrypt hash of a plaintext password, and stores both the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// Checks whether the provided plaintext password matches the hashed password stored in the struct, returning true if it matches and false otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// Add a placeholder method for inserting a new record in the table.
func (m UserModel) Insert(user *User) (err error) {
	query := `
	INSERT INTO users.users (name, surname, email, password)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`

	args := []interface{}{user.Name, user.Surname, user.Email, user.Password.hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		switch {
		default:
			return err
		}
	}
	return nil
}

// Get all the models
func (m UserModel) Index() (users []*User, err error) {
	query := `
	SELECT *
	FROM users.users
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Password.hash,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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
