package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// structure to connect to the mongo db collections
type Podcast struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// PodcastList gets all the Podcast from the db then encodes them in a response client
func (p Podcast) PodcastList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	podcastList, err := podcast.List(ctx, p.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, podcastList, http.StatusOK)
}

// Retrieve gets the Podcast from the db by _id then encodes them in a response client
func (p Podcast) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	podcastFound, err := podcast.Retrieve(ctx, p.DB, _id)
	if err != nil {
		switch err {
		case podcast.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case podcast.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for podcast %q", _id)
		}
	}

	return web.Respond(ctx, w, podcastFound, http.StatusOK)
}

// CreatePodcast decode a JSON document from a POST request and create new Podcast
// BUG: Will create empty object!!! Validate content before accepting
func (p Podcast) CreatePodcast(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var newPodcast podcast.NewPodcast

	if err := web.Decode(r, &newPodcast); err != nil {
		return err
	}

	podcast, err := podcast.CreatePodcast(ctx, p.DB, newPodcast, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, podcast, http.StatusCreated)

}

// UpdateOnePodcast decodes the body of a request to update an existing podcast.
// The ID of the podcast is part of the request URL
func (p *Podcast) UpdateOnePodcast(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	podcastID := chi.URLParam(r, "_id")

	var podcastUpdate podcast.UpdatePodcast
	if err := web.Decode(r, &podcastUpdate); err != nil {
		return errors.Wrap(err, "decoding podcast update")
	}

	if err := podcast.UpdateOnePodcast(ctx, p.DB, podcastID, podcastUpdate, time.Now()); err != nil {
		switch err {
		case podcast.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case podcast.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating podcast %q", podcastID)
		}
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}

// DeletePodcast removes a single podcast identified by an podcastID in the request URL
func (p *Podcast) DeletePodcast(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	podcastID := chi.URLParam(r, "_id")

	if err := podcast.DeletePodcast(ctx, p.DB, podcastID); err != nil {
		switch err {
		case podcast.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case podcast.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating podcast %q", podcastID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
