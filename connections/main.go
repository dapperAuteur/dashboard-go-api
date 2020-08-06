package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// connect to db
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// this formats the client
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

	// if there's an error do log the error
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database(("quickstart"))
	podcastsCollection := database.Collection("podcasts")
	episodesCollection := database.Collection("episodes")

	var podcasts []bson.M
	podcastCursor, err := podcastsCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	if err = podcastCursor.All(ctx, &podcasts); err != nil {
		panic(err)
	}
	fmt.Println(podcasts)

	var episodes []bson.M
	episodeCursor, err := episodesCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	if err = episodeCursor.All(ctx, &episodes); err != nil {
		panic(err)
	}
	fmt.Println(episodes)
}
