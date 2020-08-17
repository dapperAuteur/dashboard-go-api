package podcast

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EpisodeList gets all the Episodes for a specific Podcast from the db then encodes them in a response client
func EpisodeList(ctx context.Context, db *mongo.Collection) ([]Episode, error) {
	episodeList := []Episode{}

	episodeCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting episodeCursor. retrieving episode list")
	}

	if err = episodeCursor.All(ctx, &episodeList); err != nil {
		return nil, errors.Wrapf(err, "retrieving episode list")
	}

	return episodeList, nil
}

// PodcastEpisodeList gets all the Episodes for a specific Podcast from the db then encodes them in a response client
func PodcastEpisodeList(ctx context.Context, db *mongo.Collection, podcastID string) ([]Episode, error) {

	episodeList := []Episode{}

	// convert podcastID string from url var to podcastObjectID
	podcastObjectID, err := primitive.ObjectIDFromHex(podcastID)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	episodeCursor, err := db.Find(ctx, bson.M{"podcastID": podcastObjectID})
	if err != nil {
		return nil, errors.Wrapf(err, "getting episodeCursor. retrieving episode list")
	}

	if err = episodeCursor.All(ctx, &episodeList); err != nil {
		return nil, errors.Wrapf(err, "retrieving episode list")
	}

	return episodeList, nil
}

// RetrieveEpisode gets the first Episode in the db with the provided episodeID
func RetrieveEpisode(ctx context.Context, db *mongo.Collection, episodeID string) (*Episode, error) {

	var episode Episode

	fmt.Println(episodeID, "episodeID")

	// Check if episodeID is valid ObjectID

	episodeObjectID, err := primitive.ObjectIDFromHex(episodeID)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": episodeObjectID}).Decode(&episode); err != nil {
		return nil, apierror.ErrNotFound
	}

	return &episode, nil
}

// AddEpisode adds an Episode to a Podcast
func AddEpisode(ctx context.Context, db *mongo.Collection, newEpisode NewEpisode, podcastID string, now time.Time) (*Episode, error) {

	// convert podcastID string from url var to podcastObjectID
	podcastObjectID, err := primitive.ObjectIDFromHex(podcastID)
	if err != nil {
		return nil, errors.Wrapf(err, "converting string to ObjectID")
	}

	// put provided values into NewPodcast struct
	episode := Episode{
		PodcastID:   podcastObjectID,
		Title:       newEpisode.Title,
		Description: newEpisode.Description,
		Duration:    newEpisode.Duration,
		Spins:       newEpisode.Spins,
		Published:   newEpisode.Published,
		Tags:        newEpisode.Tags,
		CreatedAt:   now.UTC(),
		UpdatedAt:   now.UTC(),
	}

	episodeResult, err := db.InsertOne(ctx, episode)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Episode: %v", newEpisode)
	}
	fmt.Println("episodeReult : ", episodeResult)

	return &episode, nil
}
