package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Application version number
const version = "1.0.0"

// Config struct
type config struct {
	port int
	env  string
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *log.Logger
}

func main() {

	// Declare an instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: logger,
	}

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", app.healthcheckHandler)
	http.Handle("/", r)

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.

	srv := &http.Server{
		Handler:      app.routes(),
		Addr:         fmt.Sprintf(":%d", cfg.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
