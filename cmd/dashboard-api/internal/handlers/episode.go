package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// structture to connect to the mongo db Episode collection
type Episode struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// EpisodeList gets all the Episodes from the db of a specific Podcast.
// Then encodes them in a response client.
func (episode Episode) EpisodeList(w http.ResponseWriter, r *http.Request) error {

	// episodeList, err := episode.EpisodeList(r.Context(), e.DB)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// AddEpisode decodes a JSON document from a POST request and creates a new Episode for a specific Podcast.
// BUG: Will create empty object!!! Validate content before accepting
func (e Episode) AddEpisode(w http.ResponseWriter, r *http.Request) {

	// var newEpisode episode.Ne
}
