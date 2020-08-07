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

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// ==
	// Start Database
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := openDB()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// =========================================================================
	// Start API Service

	// database := client.Database(("palabras-express-api"))
	database := client.Database(("quickstart"))
	podcastsCollection := database.Collection("podcasts")

	service := Podcasts{db: podcastsCollection}

	api := http.Server{
		Addr:         "localhost:8080",
		Handler:      http.HandlerFunc(service.PodcastList),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
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

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
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
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

func openDB() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// formats the client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

	return client, err
}

// Podcast type is a group of related episodes
type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title  string             `bson:"title,omitempty" json:"title,omitempty"`
	Author string             `bson:"author,omitempty" json:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}

// structure to connect to the mongo db collections
type Podcasts struct {
	db *mongo.Collection
}

// PodcastList gets all the Podcasts from teh db then encodes them in a response client
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

// // Transaction is a line item on a balance sheet.
// type Transaction struct {
// 	Budget           string  `json:"budget,omitempty"`
// 	Currency         string  `json:"currency,omitempty"`
// 	FinancialAccount string  `json:"financial_account,omitempty"`
// 	Media            string  `json:"media,omitempty"`
// 	Note             string  `json:"note,omitempty"`
// 	Occurrence       string  `json:"occurrence,omitempty"`
// 	Participant      string  `json:"participant,omitempty"`
// 	Tag              string  `json:"tag,omitempty"`
// 	TransactionEvent string  `json:"transaction_event,omitempty"`
// 	TransactionValue float64 `json:"transaction_value,omitempty"`
// 	Vendor           string  `json:"vendor,omitempty"`
// }

// // Verbo is a Spanish verb
// type Verbo struct {
// 	CambiarDeIrregular   string  `json:"cambiar_de_irregular,omitempty"`
// 	CategoriaDeIrregular string  `json:"categoria_de_irregular,omitempty"`
// 	English              string  `json:"english,omitempty"`
// 	Grupo                float64 `json:"grupo,omitempty"`
// 	Irregular            bool    `json:"irregular,omitempty"`
// 	Media                string  `json:"media,omitempty"`
// 	Note                 string  `json:"note,,omitempty"`
// 	Reflexive            bool    `json:"reflexive,omitempty"`
// 	Spanish              string  `json:"spanish,omitempty"`
// 	Terminacion          string  `json:"terminacion,omitempty"`
// }

// // ListTransactions is an HTTP Handler for returning a list of Transactions.
// func ListTransactions(w http.ResponseWriter, r *http.Request) {
// 	list := []Transaction{
// 		{Budget: "Food", TransactionValue: 14.39, Vendor: "Fry's"},
// 		{Budget: "Tools", TransactionValue: 1400.39, Vendor: "System76"},
// 	}

// 	data, err := json.Marshal(list)
// 	if err != nil {
// 		log.Println("error marshalling result", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	if _, err := w.Write(data); err != nil {
// 		log.Println("error writing result", err)
// 	}
// }

// func ListVerbos(w http.ResponseWriter, r *http.Request) {
// 	verbos := []Verbo{}

// 	data, err := json.Marshal(verbos)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		log.Println("error marshalling", err)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	if _, err := w.Write(data); err != nil {
// 		log.Println("error writing", err)
// 	}
// }
