package data

import (
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Token string      `json:"string"`
	User  interface{} `json:"user"`
}

type AuthCredentials struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
}

// Define a struct type which wraps a sql.DB connection pool.
type AuthModel struct {
}

type password struct {
	plaintext *string
	hash      []byte
}

// Calculates the bcrypt hash of a plaintext password, and stores both the hash and the plaintext versions in the struct.
func Hash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
