package main

import (
	"database/sql"
	"log"
	"net/http"
	"playground/internal"
	"playground/internal/database"
	"playground/migrations"
)

func main() {
	server := newServer()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func newServer() *http.Server {
	db, err := initialiseDB()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on port %s", ":5000")
	log.Printf("listening on port %s", ":5000")

	router := http.NewServeMux()
	internal.ConfigureRoutes(router, db)

	return &http.Server{
		Handler: router,
		Addr:    ":5000",
	}
}

func initialiseDB() (*sql.DB, error) {
	db, err := database.ConnectDB()

	if err != nil {
		return nil, err
	}

	if err = migrations.MigrateUp(db); err != nil {
		return nil, err
	}

	return db, nil
}
