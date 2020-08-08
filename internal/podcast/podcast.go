package podcast

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
