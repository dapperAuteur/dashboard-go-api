package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Podcast structure to connect to the mongo db collections
type Podcast struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// PodcastList gets all the Podcast from the db then encodes them in a response client
func (p Podcast) PodcastList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Podcast.PodcastList")
	defer span.End()

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

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var newPodcast podcast.NewPodcast

	if err := web.Decode(r, &newPodcast); err != nil {
		return err
	}

	podcast, err := podcast.CreatePodcast(ctx, p.DB, claims, newPodcast, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, podcast, http.StatusCreated)

}

// UpdateOnePodcast decodes the body of a request to update an existing podcast.
// The ID of the podcast is part of the request URL
func (p *Podcast) UpdateOnePodcast(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	fmt.Print("*****    UpdateOnePodcast   *****\n")

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}
	fmt.Print("*****    claims    *****\n")

	podcastID := chi.URLParam(r, "_id")

	var podcastUpdate podcast.UpdatePodcast
	if err := web.Decode(r, &podcastUpdate); err != nil {
		return errors.Wrap(err, "decoding podcast update")
	}

	if err := podcast.UpdateOnePodcast(ctx, p.DB, claims, podcastID, podcastUpdate, time.Now()); err != nil {
		switch err {
		case podcast.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case podcast.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case podcast.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
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
