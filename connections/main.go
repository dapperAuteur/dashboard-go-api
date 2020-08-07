package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// connect to db
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// // this formats the client
	// // client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))

	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

	client, err := openDB()

	// if there's an error do log the error
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database(("quickstart"))
	podcastsCollection := database.Collection("podcasts")
	// episodesCollection := database.Collection("episodes")

	service := Podcasts{db: podcastsCollection}

	// ==
	// Start API Service

	api := http.Server{
		Addr:         "localhost:8080",
		Handler:      http.HandlerFunc(service.PodcastList),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// ==
	// Shutdown

	// Blocking main and waiting for shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not conpmlete in %v : %v", timeout, err)
			err = api.Close()
		}
	}

	// var podcasts []Podcast
	// podcastCursor, err := podcastsCollection.Find(ctx, bson.M{})
	// if err != nil {
	// 	panic(err)
	// }

	// if err = podcastCursor.All(ctx, &podcasts); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(podcasts)

	// var episodes []Episode
	// episodeCursor, err := episodesCollection.Find(ctx, bson.M{
	// 	"duration": bson.D{{"$gt", 25}},
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// if err = episodeCursor.All(ctx, &episodes); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(episodes)

	// Podcasts holds business logic related to Podcasts

}

func openDB() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// this formats the client
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

	return client, err
}

// Podcast type is a group of related episodes
type Podcast struct {
	ID     primitive.ObjectID `bson: "_id, omitempty" json:"_id, omitempty"`
	Title  string             `bson: "title, omitempty" json: "title, omitempty"`
	Author string             `bson: "author, omitempty" json: "author, omitempty"`
	Tags   []string           `bson: "tags, omitempty" json: "tags, omitempty"`
}

// Episode is the video or audio content
type Episode struct {
	ID          primitive.ObjectID `bson: "_id, omitempty" json: "_id, omitempty"`
	Podcast     primitive.ObjectID `bson: "podcast, omitempty" json: "podcast, omitempty"`
	Title       string             `bson: "title, omitempty" json: "title, omitempty"`
	Description string             `bson: "description, omitempty" json: "description, omitempty"`
	Duration    int32              `bson: "duration, omitempty" json: "duration, omitempty"`
}

type Podcasts struct {
	db *mongo.Collection
}

// PodcastList gets all the Podcasts from the db then encodes them in a response client
func (p Podcasts) PodcastList(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	podcastList := []Podcast{}

	podcastCursor, err := p.db.Find(ctx, bson.M{})
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
