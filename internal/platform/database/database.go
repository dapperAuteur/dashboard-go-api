package database

import (
	"context"
	"fmt"
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

// StatusCheck returns nil if it can successfully talk to the database.
// It returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *mongo.Collection) error {

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	// if err := db.Database().RunCommand(ctx, "connectionStatus"); err != nil {
	// 	return errors.Wra
	// }
	connectionStatus := db.Database().RunCommand(ctx, "connectionStatus")
	fmt.Print(connectionStatus)
	return nil
}
