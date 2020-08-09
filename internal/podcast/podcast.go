package podcast

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/mongo"
)

// PodcastList gets all the Podcasts from the db then encodes them in a response client
func List(db *mongo.Collection) ([]Podcast, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	podcastList := []Podcast{}

	podcastCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	if err = podcastCursor.All(ctx, &podcastList); err != nil {
		return nil, err
	}

	return podcastList, nil
}

// Retrieve gets the first Podcast in the db with the provided _id
func Retrieve(db *mongo.Collection, _id string) (*Podcast, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var podcast Podcast

	id, err := primitive.ObjectIDFromHex(_id)
	err := db.FindOne(ctx, podcastFound).Decode(&podcast)
	if err != nil {
		return nil, err
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&podcast); err != nil {
		log.Printf("podcast not found: %s", podcast)
		log.Printf("id sent to podcast.Retrieve podcast}: %s", podcast)
		return nil, err
	}

	fmt.Println("result AFTER:", podcast)

	return &podcast, nil
}

	return &podcast, nil
}
