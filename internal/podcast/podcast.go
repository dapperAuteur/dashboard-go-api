package podcast

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/mongo"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("product not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// List gets all the Podcasts from the db then encodes them in a response client
func List(ctx context.Context, db *mongo.Collection) ([]Podcast, error) {
	podcastList := []Podcast{}

	podcastCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting podcastCursor. retrieving podcast list")
	}

	if err = podcastCursor.All(ctx, &podcastList); err != nil {
		return nil, errors.Wrapf(err, "retrieving podcast list")
	}

	return podcastList, nil
}

// Retrieve gets the first Podcast in the db with the provided _id
func Retrieve(ctx context.Context, db *mongo.Collection, _id string) (*Podcast, error) {

	var podcast Podcast

	// Check if _id is valid ObjectID
	// if _, err := uuid.Parse(_id); err != nil {
	// 	return nil, ErrInvalidID
	// }

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&podcast); err != nil {
		log.Printf("podcast not found: %s", podcast)
		log.Printf("id sent to podcast.Retrieve podcast}: %s", id)
		return nil, ErrNotFound
	}

	fmt.Println("result AFTER:", podcast)

	return &podcast, nil
}

// RetrieveByTitle gets the first Podcast in the db with the provided title
func RetrieveByTitle(ctx context.Context, db *mongo.Collection, title string) (*Podcast, error) {

	var podcast Podcast

	// This works to find the Podcast by title
	filter := Podcast{
		Title: "The Nic Raboy Show",
	}

	if err := db.FindOne(ctx, filter).Decode(&podcast); err != nil {
		log.Printf("podcast title not found: %s", title)
		log.Printf("id sent to podcast.Retrieve podcast}: %s", podcast)
		return nil, errors.Wrapf(err, "retrieving podcast by title: %s", title)
	}

	fmt.Println("result AFTER:", podcast)

	return &podcast, nil
}

// CreatePodcast will create a new Podcast in the database and returns the new Podcast
func CreatePodcast(ctx context.Context, db *mongo.Collection, newPodcast NewPodcast, now time.Time) (*Podcast, error) {

	podcast := Podcast{
		Title:       newPodcast.Title,
		Author:      newPodcast.Author,
		Subscribers: newPodcast.Subscribers,
		Tags:        newPodcast.Tags,
		Published:   newPodcast.Published,
		CreatedAt:   now.UTC(),
		UpdatedAt:   now.UTC(),
	}

	// How do I get MongoDB to return the new Podcast
	podcastResult, err := db.InsertOne(ctx, podcast)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Podcast: %v", newPodcast)
	}
	fmt.Println("podcastResult : ", podcastResult)

	// doesn't return ObjectID with podcast, find a way to get the _id with the Podcast
	return &podcast, nil
}
