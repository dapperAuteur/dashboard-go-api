package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Podcasts defines all of the handlers related to podcasts. It holds the
// application state needed by the handler methods.
type Podcast struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// PodcastList gets all the Podcast from the service layer.
func (p Podcast) PodcastList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// panic("OH NO!!!") // create an artificial PANIC

	ctx, span := trace.StartSpan(ctx, "handlers.Podcast.PodcastList")
	defer span.End()

	podcastList, err := podcast.List(ctx, p.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, podcastList, http.StatusOK)
}

// Retrieve gets the Podcast from the db identified by an _id in the request URL, then encodes it in a response client
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

// RetrieveByTitle gets the Podcast from the db identified by an title in the request URL, then encodes it in a response client
func (p Podcast) RetrieveByTitle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	title := chi.URLParam(r, "title")

	podcastFound, err := podcast.Retrieve(ctx, p.DB, title)
	if err != nil {
		switch err {
		case podcast.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case podcast.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for podcast %q", title)
		}
	}

	return web.Respond(ctx, w, podcastFound, http.StatusOK)
}

// CreatePodcast decodes the body of a request to create a new podcast.
// The full podcast with generated fields is sent back in the response.
func (p Podcast) CreatePodcast(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
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

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	podcastID := chi.URLParam(r, "_id")

	if err := podcast.DeletePodcast(ctx, p.DB, claims, podcastID); err != nil {
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

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
