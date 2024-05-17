package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMqEnv struct {
	mqUser string
	mqPass string
	mqHost string
	mqPort string
}

type Auth struct {
	Token    string      `json:"string"`
	User     interface{} `json:"user"`
	Datetime time.Time   `json:"datetime"`
}

type application struct {
	logger      *log.Logger
	rabbitMqEnv *rabbitMqEnv
}

func main() {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		logger: logger,
		rabbitMqEnv: &rabbitMqEnv{
			mqUser: os.Getenv("RABBITMQ_USER"),
			mqPass: os.Getenv("RABBITMQ_PASS"),
			mqHost: os.Getenv("RABBITMQ_HOST"),
			mqPort: os.Getenv("RABBITMQ_PORT"),
		},
	}

	logger.Printf("Notifications microservice started")
	logger.Printf("--- RabbitMQ Configuration: ---")
	logger.Printf("Host: %s", app.rabbitMqEnv.mqHost)
	logger.Printf("Port: %s", app.rabbitMqEnv.mqPort)
	logger.Printf("User: %s", app.rabbitMqEnv.mqUser)
	logger.Printf("---------------------")

	conn, err := initMQConnection(app.rabbitMqEnv)
	if err != nil {
		app.logger.Printf("Failed to connect to RabbitMQ: %s", err)
	}

	defer conn.Close()

	ch, err := initMQChannel(conn)
	if err != nil {
		app.logger.Printf("Failed to open a channel: %s", err)
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		"user_login_event", // name
		"fanout",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)

	if err != nil {
		app.logger.Fatalf("Failed to declare an exchange: %s", err)
	}

	q, err := ch.QueueDeclare(
		"login_email_notification", // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)

	if err != nil {
		app.logger.Fatalf("Failed to declare a queue: %s", err)
	}

	err = ch.QueueBind(
		q.Name,             // queue name
		"",                 // routing key
		"user_login_event", // exchange
		false,
		nil)

	if err != nil {
		app.logger.Fatalf("Failed to bind a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		app.logger.Fatalf("Failed to register a consumer: %s", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			err := sendLoginEmailNotification(d.Body)
			if err != nil {
				d.Nack(false, true)
				app.logger.Fatalf("Failed to handle queue message: %s", err)
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

// initMQChannel initializes a channel to RabbitMQ
func initMQChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// initMQConnection initializes a connection to RabbitMQ
func initMQConnection(env *rabbitMqEnv) (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://" + env.mqUser + ":" + env.mqPass + "@" + env.mqHost + ":" + env.mqPort + "/")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func sendLoginEmailNotification(body []byte) error {
	auth := &Auth{}
	err := json.Unmarshal(body, &auth)
	if err != nil {
		return err
	}

	// Sleep random time between 1 and 5 seconds
	time.Sleep(time.Duration(rand.Intn(5-1)+1) * time.Second)
	log.Printf("Email sent to: %s", auth.User)

	return nil
}
