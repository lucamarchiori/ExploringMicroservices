package data

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

func (auth *Auth) TriggerLoginEvent() error {
	log.Println("Triggering login event")

	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@" + host + ":" + port + "/")
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Connected to RabbitMQ at %s:%s", host, port)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	log.Println("Opened a channel")

	err = ch.ExchangeDeclare(
		"user_events", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		return err
	}

	log.Println("Declared an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serialized, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"user_events", // exchange
		"login",       // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        serialized,
		})

	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", serialized)

	return nil
}
