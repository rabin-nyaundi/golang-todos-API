package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4010, "API Server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	// mux := http.NewServeMux()
	// mux.HandleFunc("/v1/todos", app.todo)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("application server started")

	logger.Printf("starting %s server on %d", app.config.env, app.config.port)

	err := srv.ListenAndServe()
	logger.Fatal(err)
	fmt.Println("Started your todo application")
}
