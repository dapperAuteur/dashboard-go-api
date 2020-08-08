package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// ==
	// Start Database
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := openDB()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

}

// Do we have a use for a switch statement to run admin commands?

func openDB() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// formats the client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

	return client, err
}
