package main

import (
	"auth-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const port = "80"

var counts int64

type Config struct {
	Repo data.Repository
}

func main() {
	log.Println("Starting auth-service")

	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Panic(err)
	}

	// Setup config
	app := Config{}
	app.setupRepo(db)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	fmt.Println(dsn)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Postgres ready!")
			return connection, nil
		}

		if counts > 5 {
			log.Fatal("Too many Postgres connection attempts")
			return nil, err
		}

		log.Println("Retrying in 5 second...")
		time.Sleep(5 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(db *sql.DB) {
	repo := data.New(db)
	app.Repo = repo
}
