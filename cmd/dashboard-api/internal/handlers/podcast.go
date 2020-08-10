package handlers

import (
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
func (p Podcast) PodcastList(w http.ResponseWriter, r *http.Request) error {

	podcastList, err := podcast.List(r.Context(), p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, podcastList, http.StatusOK)
}

// Retrieve gets the Podcast from the db by _id then encodes them in a response client
func (p Podcast) Retrieve(w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	podcastFound, err := podcast.Retrieve(r.Context(), p.DB, _id)
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

	return web.Respond(w, podcastFound, http.StatusOK)
}

// CreatePodcast decode a JSON document from a POST request and create new Podcast
// BUG: Will create empty object!!! Validate content before accepting
func (p Podcast) CreatePodcast(w http.ResponseWriter, r *http.Request) error {

	var newPodcast podcast.NewPodcast

	if err := web.Decode(r, &newPodcast); err != nil {
		return err
	}

	podcast, err := podcast.CreatePodcast(r.Context(), p.DB, newPodcast, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, podcast, http.StatusCreated)

}
