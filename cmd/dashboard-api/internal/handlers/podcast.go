package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// structure to connect to the mongo db collections
type Podcast struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// PodcastList gets all the Podcast from the db then encodes them in a response client
func (p Podcast) PodcastList(w http.ResponseWriter, r *http.Request) {
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
		p.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing result", err)
	}
}

// Retrieve gets the Podcast from the db by _id then encodes them in a response client
func (p Podcast) Retrieve(w http.ResponseWriter, r *http.Request) {
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_id := chi.URLParam(r, "_id")

	podcast, err := podcast.Retrieve(p.DB, _id)
	if err != nil {
		panic(err)
	}

	fmt.Println(podcast)

	data, err := json.Marshal(podcast)
	if err != nil {
		p.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing result", err)
	}
}

// CreatePodcast decode a JSON document from a POST request and create new Podcast
func (p Podcast) CreatePodcast(w http.ResponseWriter, r *http.Request) {

	var newPodcast podcast.NewPodcast

	if err := json.NewDecoder(r.Body).Decode(&newPodcast); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	podcast, err := podcast.CreatePodcast(p.DB, newPodcast, time.Now())
	if err != nil {
		panic(err)
	}

	fmt.Println(podcast)

	data, err := json.Marshal(podcast)
	if err != nil {
		p.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing result", err)
	}

}
