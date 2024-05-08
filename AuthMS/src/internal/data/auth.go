package data

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

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
func HashPassword(plaintextPassword string) (password string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}
	password = string(hash)
	return password, nil
}

// Checks whether the provided plaintext password matches the hashed password stored in the struct, returning true if it matches and false otherwise.
func PasswordMatch(plaintextPassword string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintextPassword))
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

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
