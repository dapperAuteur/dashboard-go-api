package handlers

import (
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/mid"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
)

// API constructs a handler that knows about all API routes.
func API(logger *log.Logger, db *mongo.Database, authenticator *auth.Authenticator) http.Handler {

	app := web.NewApp(logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics())

	c := Check{DB: db.Collection("podcasts")}

	app.Handle(http.MethodGet, "/v1/health", c.Health)

	u := Users{DB: db.Collection("users"), authenticator: authenticator}

	app.Handle(http.MethodGet, "/v1/users/token", u.Token)

	// episodesCollection := db.Collection("episodes")
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

	app.Handle(http.MethodGet, "/v1/episodes", episode.EpisodeList)
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}/episodes", episode.PodcastEpisodeList)
	app.Handle(http.MethodGet, "/v1/episodes/{episodeID}", episode.RetrieveEpisode)
	app.Handle(http.MethodPost, "/v1/podcasts/{_id}/episodes", episode.AddEpisode)
	// app.Handle(http.MethodGet, "/v1/podcasts/{_id}/episodes/{_id}", episode.Retrieve)

	app.Handle(http.MethodGet, "/v1/podcasts", podcast.PodcastList)
	app.Handle(http.MethodPost, "/v1/podcasts", podcast.CreatePodcast)
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}", podcast.Retrieve)
	app.Handle(http.MethodPut, "/v1/podcasts/{_id}", podcast.UpdateOnePodcast)
	app.Handle(http.MethodDelete, "/v1/podcasts/{_id}", podcast.DeletePodcast)

	return app
}
