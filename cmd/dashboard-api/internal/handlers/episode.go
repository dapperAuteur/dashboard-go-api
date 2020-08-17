package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// Episode structture to connect to the mongo db Episode collection
type Episode struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// EpisodeList gets all the Episodes from the db of all Podcasts.
// Then encodes them in a response client.
func (e Episode) EpisodeList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	episodeList, err := podcast.EpisodeList(r.Context(), e.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, episodeList, http.StatusOK)
}

// PodcastEpisodeList gets all the Episodes from the db of a specific Podcast.
// Then encodes them in a response client.
func (e Episode) PodcastEpisodeList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	podcastID := chi.URLParam(r, "_id")

	episodeList, err := podcast.PodcastEpisodeList(ctx, e.DB, podcastID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, episodeList, http.StatusOK)
}

// RetrieveEpisode gets the Episode from the db by episodeID then encodes it in a response client
func (e Episode) RetrieveEpisode(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	episodeID := chi.URLParam(r, "episodeID")

	episodeFound, err := podcast.RetrieveEpisode(ctx, e.DB, episodeID)
	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for podcast episode %q", episodeID)
		}
	}
	return web.Respond(ctx, w, episodeFound, http.StatusOK)
}

// AddEpisode decodes a JSON document from a POST request and creates a new Episode for a specific Podcast.
// BUG: Will create empty object!!! Validate content before accepting
func (e Episode) AddEpisode(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var newEpisode podcast.NewEpisode

	podcastID := chi.URLParam(r, "_id")

	if err := web.Decode(r, &newEpisode); err != nil {
		return err
	}

	episode, err := podcast.AddEpisode(ctx, e.DB, newEpisode, podcastID, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, episode, http.StatusCreated)
}
