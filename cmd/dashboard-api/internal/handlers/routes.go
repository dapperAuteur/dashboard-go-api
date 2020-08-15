package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/dapperAuteur/dashboard-go-api/internal/mid"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
)

// API constructs a handler that knows about all API routes.
func API(shutdown chan os.Signal, logger *log.Logger, db *mongo.Database, authenticator *auth.Authenticator) http.Handler {

	app := web.NewApp(shutdown, logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics(), mid.Panics(logger))

	c := Check{DB: db.Collection("podcasts")}

	app.Handle(http.MethodGet, "/v1/health", c.Health)

	u := Users{DB: db.Collection("users"), authenticator: authenticator}

	app.Handle(http.MethodGet, "/v1/users/token", u.Token)

	// episodesCollection := db.Collection("episodes")
	budgetsCollection := db.Collection("budgets")
	episodesCollection := db.Collection("episodes")
	podcastsCollection := db.Collection("podcasts")

	budget := Budget{
		DB:  budgetsCollection,
		Log: logger,
	}

	podcast := Podcast{
		DB:  podcastsCollection,
		Log: logger,
	}

	episode := Episode{
		DB:  episodesCollection,
		Log: logger,
	}

	// Budget Routes
	app.Handle(http.MethodGet, "/v1/budgets", budget.List)
	app.Handle(http.MethodPost, "/v1/budgets", budget.Create, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// Episode Routes
	app.Handle(http.MethodGet, "/v1/episodes", episode.EpisodeList, mid.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}/episodes", episode.PodcastEpisodeList, mid.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/episodes/{episodeID}", episode.RetrieveEpisode, mid.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/podcasts/{_id}/episodes", episode.AddEpisode, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	// app.Handle(http.MethodGet, "/v1/podcasts/{_id}/episodes/{_id}", episode.Retrieve, mid.Authenticate(authenticator))

	// Podcast Routes
	app.Handle(http.MethodGet, "/v1/podcasts", podcast.PodcastList, mid.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/podcasts", podcast.CreatePodcast, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}", podcast.Retrieve, mid.Authenticate(authenticator))
	app.Handle(http.MethodPut, "/v1/podcasts/{_id}", podcast.UpdateOnePodcast, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/v1/podcasts/{_id}", podcast.DeletePodcast, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	return app
}
