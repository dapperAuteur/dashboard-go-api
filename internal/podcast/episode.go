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

	podcastObjectID, err := primitive.ObjectIDFromHex(podcastID)
	if err != nil {
		return nil, errors.Wrapf(err, "converting string to ObjectID")
	}

	episode := NewEpisode{
		PodcastID:   podcastObjectID,
		Title:       newEpisode.Title,
		Description: newEpisode.Description,
		Duration:    newEpisode.Duration,
		Tags:        newEpisode.Tags,
	}

	fmt.Println("newEpisode : " epispode),
	// use db connection to insertOne
	return &episode
}
