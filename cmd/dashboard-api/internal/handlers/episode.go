package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

// structture to connect to the mongo db Episode collection
type Episode struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// EpisodeList gets all the Episodes from the db of a specific Podcast.
// Then encodes them in a response client.
func (e Episode) EpisodeList(w http.ResponseWriter, r *http.Request) error {

	// episodeList, err := episode.EpisodeList(r.Context(), e.DB)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// AddEpisode decodes a JSON document from a POST request and creates a new Episode for a specific Podcast.
// BUG: Will create empty object!!! Validate content before accepting
func (e Episode) AddEpisode(w http.ResponseWriter, r *http.Request) error {

	var newEpisode podcast.NewEpisode

	podcastID := chi.URLParam(r, "_id")

	if err := web.Decode(r, &newEpisode); err != nil {
		return err
	}

	episode, err := podcast.AddEpisode(r.Context(), e.DB, newEpisode, podcastID, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, episode, http.StatusCreated)
}
