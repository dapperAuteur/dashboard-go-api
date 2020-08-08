package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config is what's required to open database connection
type Config struct {
	AtlasURI string
}

// Open knows how to open a database connection
func Open(cfg Config) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// formats the client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.AtlasURI))

	return client, err
}
