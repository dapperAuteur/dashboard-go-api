package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Podcast struct {
	ID     primitive.ObjectID `bson: "_id, omitempty"`
	Title  string             `bson: "title, omitempty"`
	Author string             `bson: "author, omitempty"`
	Tags   []string           `bson: "tags, omitempty"`
}

type Episode struct {
	ID          primitive.ObjectID `bson: "_id, omitempty"`
	Podcast     primitive.ObjectID `bson: "podcast, omitempty"`
	Title       string             `bson: "title, omitempty"`
	Description string             `bson: "description, omitempty"`
	Duration    int32              `bson: "duration, omitempty"`
}

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

	var podcasts []Podcast
	podcastCursor, err := podcastsCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	if err = podcastCursor.All(ctx, &podcasts); err != nil {
		panic(err)
	}
	fmt.Println(podcasts)

	// create Podcast
	podcast := Podcast{
		Title:  "The Polyglot Developer",
		Author: "Nic Raboy",
		Tags:   []string{"development", "programming", "coding"},
	}
	// insert new podcast into db
	insertResult, err := podcastsCollection.InsertOne(ctx, podcast)
	if err != nil {
		panic(err)
	}
	// print id of inserted id
	fmt.Println(insertResult.InsertedID)

	var episodes []Episode
	episodeCursor, err := episodesCollection.Find(ctx, bson.M{
		"duration": bson.D{{"$gt", 25}},
	})
	if err != nil {
		panic(err)
	}
	if err = episodeCursor.All(ctx, &episodes); err != nil {
		panic(err)
	}
	fmt.Println(episodes)
}
