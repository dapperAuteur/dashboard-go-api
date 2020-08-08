package handlers

import (
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
)

// API constructs a handler that knows about all API routes.
func API(logger *log.Logger, db *mongo.Collection) http.Handler {

	app := web.NewApp(logger)

	podcast := Podcast{
		DB:  db,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/v1/podasts", podcast.PodcastList)
	app.Handle(http.MethodGet, "/v1/podasts/{_id}", podcast.Retrieve)

	return app
}
