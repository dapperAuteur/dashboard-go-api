package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// structure to connect to the mongo db collections
type Podcasts struct {
	DB *mongo.Collection
}

// PodcastList gets all the Podcasts from the db then encodes them in a response client
func (p Podcasts) PodcastList(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	podcastList, err := podcast.List(p.DB)
	if err != nil {
		panic(err)
	}

	podcastCursor, err := p.DB.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	if err = podcastCursor.All(ctx, &podcastList); err != nil {
		panic(err)
	}
	fmt.Println(podcastList)

	data, err := json.Marshal(podcastList)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}
