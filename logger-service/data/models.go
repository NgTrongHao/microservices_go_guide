package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func New(mongoClient *mongo.Client) Models {
	client = mongoClient

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) GetCollection() *mongo.Collection {
	return client.Database("logs").Collection("logs")
}

func (l *LogEntry) Insert(logEntry LogEntry) error {
	collection := l.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Error inserting log entry: %v", err)
		return err
	}
	return nil
}

func (l *LogEntry) GetAll() ([]LogEntry, error) {
	collection := l.GetCollection()

	// sort by created_at descending
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.Background(), struct{}{}, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(cursor, context.Background())

	var logs []LogEntry
	for cursor.Next(context.Background()) {
		var entry LogEntry
		err := cursor.Decode(&entry)
		if err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}
	return logs, nil
}

func (l *LogEntry) GetByID(id string) (LogEntry, error) {
	collection := l.GetCollection()
	var entry LogEntry
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&entry)
	return entry, err
}

func (l *LogEntry) DeleteAll() error {
	collection := l.GetCollection()
	_, err := collection.DeleteMany(context.Background(), struct{}{})
	return err
}
