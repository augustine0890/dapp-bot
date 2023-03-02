package database

import (
	"context"
	"fmt"
	"time"

	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMongoClient returns a new MongoDB client instance
func GetMongoClient(uri string, dbName string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logging.Info("Connected to MongoDB")

	return client, nil
}

// GetUsersColl returns the MongoDB collection of users
func GetUsersColl(mongoClient *mongo.Client, cfg *config.Config) *mongo.Collection {
	return mongoClient.Database(cfg.MongoDBName).Collection("users")
}
