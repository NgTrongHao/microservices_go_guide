package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	port      = "80"
	rpc_port  = "5001"
	grpc_port = "50001"
)

var (
	db_user = os.Getenv("DB_USER")
	db_pass = os.Getenv("DB_PASS")
	db_url  = os.Getenv("DB_URL")
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoClient, err := connectToMongoDB()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Panic(err)
		}
	}(client, ctx)

	app := Config{
		Models: data.New(mongoClient),
	}

	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db_url)
	clientOptions.SetAuth(options.Credential{
		Username: db_user,
		Password: db_pass,
	})
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
