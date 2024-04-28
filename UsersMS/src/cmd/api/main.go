package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"lucamarchiori/MicroserviceBoilerplate/internal/data"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Application version number
const version = "1.0.0"

// Config struct
type config struct {
	port int
	env  string
	db   struct {
		host     string
		port     int
		user     string
		pass     string
		database string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	// Declare an instance of the config struct.
	var cfg config
	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Parse environment variables for database connection
	if os.Getenv("POSTGRES_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		if err != nil {
			log.Fatal(err)
		}
		flag.IntVar(&cfg.db.port, "db-port", port, "Database port")
	}
	flag.StringVar(&cfg.db.user, "db-user", os.Getenv("POSTGRES_USER"), "Database user")
	flag.StringVar(&cfg.db.pass, "db-pass", os.Getenv("POSTGRES_PASSWORD"), "Database password")
	flag.StringVar(&cfg.db.database, "db-database", os.Getenv("POSTGRES_DB"), "Database name")
	flag.StringVar(&cfg.db.host, "db-host", os.Getenv("POSTGRES_HOST"), "Database host")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	logger.Printf("UsersMS version %s started", version)
	logger.Printf("--- Configuration: ---")
	logger.Printf("Server port: %d", cfg.port)
	logger.Printf("Environment: %s", cfg.env)
	logger.Printf("Database host: %s", cfg.db.host)
	logger.Printf("Database port: %d", cfg.db.port)
	logger.Printf("Database user: %s", cfg.db.user)
	logger.Printf("Database name: %s", cfg.db.database)
	logger.Printf("---------------------")

	logger.Printf("Starting database connection procedure ...")

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool established")

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
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
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.db.host, cfg.db.port, cfg.db.user, cfg.db.pass, cfg.db.database))
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
