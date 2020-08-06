package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// this formats the client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

	// if there's an error do log the error
	if err != nil {
		log.Fatal(err)
	}

	// connect to db
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database(("quickstart"))
	podcastsCollection := database.Collection("podcasts")
	// 5f2b44a5f15bb82f27e4c0f3

	// convert the id to something mongo go driver understand
	id, _ := primitive.ObjectIDFromHex("5f2b44a5f15bb82f27e4c0f3")

	result, err := podcastsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"author", "Nicolas Raboy"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	result, err = podcastsCollection.UpdateMany(
		ctx,
		bson.M{"title": "The Polyglot Developer Podcast"},
		bson.D{
			{"$set", bson.D{{"author", "Nic Raboy"}}},
		},
	)
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

}
