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
	Password  string    `json:"password,omitempty"`
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
func HashPassword(plaintextPassword string) (password string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}
	password = string(hash)
	return password, nil
}

// Checks whether the provided plaintext password matches the hashed password stored in the struct, returning true if it matches and false otherwise.
func (u User) passwordMatch(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plaintextPassword))
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

	args := []interface{}{user.Name, user.Surname, user.Email, user.Password}
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
func (m UserModel) Index(email string) (users []*User, err error) {
	query := `
	SELECT *
	FROM users.users
	WHERE (LOWER(email) = LOWER($1) OR $1 = '')
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Password,
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

func (m UserModel) Get(id int64) (*User, error) {
	query := `
	SELECT *
	FROM users.users
	WHERE id=$1
	`
	var u User
	err := m.DB.QueryRow(query, id).Scan(
		&u.Id,
		&u.Name,
		&u.Surname,
		&u.Email,
		&u.Password,
		&u.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}

// Add a placeholder method for updating a specific record in the table.
func (m UserModel) Update(user *User) error {
	// TBI
	return nil
}

// Add a placeholder method for deleting a specific record from the table.
func (m UserModel) Delete(id int64) error {
	// TBI
	return nil
}
