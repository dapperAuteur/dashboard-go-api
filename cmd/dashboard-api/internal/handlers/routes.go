package handlers

import (
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
)

// API constructs a handler that knows about all API routes.
func API(logger *log.Logger, db *mongo.Database) http.Handler {

	app := web.NewApp(logger)

	episodesCollection := db.Collection("episodes")
	podcastsCollection := db.Collection("podcasts")

	podcast := Podcast{
		DB:  podcastsCollection,
		Log: logger,
	}

	episode := Episode{
		DB:  episodesCollection,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/v1/podcasts/{_id}/v1/episodes", episode.EpisodeList)
	// app.Handle(http.MethodPost, "/v1/podcasts/{_id}/v1/episodes", episode.AddEpisode)
	// app.Handle(http.MethodGet, "/v1/podcasts/{_id}/v1/episodes/{_id}", episode.Retrieve)

	app.Handle(http.MethodGet, "/v1/podcasts", podcast.PodcastList)
	app.Handle(http.MethodPost, "/v1/podcasts", podcast.CreatePodcast)
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}", podcast.Retrieve)

	return app
}
