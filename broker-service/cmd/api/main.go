package main

import (
	"log"
	"net/http"
)

const port = "80"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker-service on port %s\n", port)

	// Define http server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
