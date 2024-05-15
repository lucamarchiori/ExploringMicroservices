package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMqEnv struct {
	mqUser string
	mqPass string
	mqHost string
	mqPort string
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

	ch, err := initMQConnection(app.rabbitMqEnv)
	if err != nil {
		app.logger.Printf("Failed to connect to RabbitMQ: %s", err)
	}

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
		app.logger.Fatalf("Failed to declare an exchange: %s", err)
	}

	q, err := ch.QueueDeclare(
		"email_notification", // name
		false,                // durable
		false,                // delete when unused
		true,                 // exclusive
		false,                // no-wait
		nil,                  // arguments
	)

	if err != nil {
		app.logger.Fatalf("Failed to declare a queue: %s", err)
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		"login",       // routing key
		"user_events", // exchange
		false,
		nil)

	if err != nil {
		app.logger.Fatalf("Failed to bind a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
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
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func initMQConnection(env *rabbitMqEnv) (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://" + env.mqUser + ":" + env.mqPass + "@" + env.mqHost + ":" + env.mqPort + "/")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()

	return ch, nil
}
