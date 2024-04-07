package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
	clientErr  error
	dbCol      *mongo.Collection
)

func getClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		client, clientErr = connectToMongoDB()
	})

	// If clientErr is not nil (indicating an error occurred previously), attempt to reconnect
	if clientErr != nil {
		client, clientErr = connectToMongoDB()
	}

	return client, clientErr
}

func connectToMongoDB() (*mongo.Client, error) {
	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(100)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Ping the MongoDB server
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Access a database
	database := client.Database(os.Getenv("DB_NAME"))

	// Access a collection
	dbCol = database.Collection(os.Getenv("COLLECTION_NAME"))

	return client, nil
}
