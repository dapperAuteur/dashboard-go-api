package podcast

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
