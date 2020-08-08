package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/cmd/dashboard-api/internal/handlers"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/conf"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/database"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := log.New(os.Stdout, "DASHBOARD : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// ==
	// Configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			AtlasUri string `conf:"default:mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"`
		}
	}

	// ==
	// Get Configuration
	// Helpful info in case of error

	if err := conf.Parse(os.Args[1:], "DASHBOARD", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DASHBOARD", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// print config values when app starts
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// is it ok to do this twice
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := database.Open(database.Config{
		AtlasUri: cfg.DB.AtlasUri,
	})
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// =========================================================================
	// Start API Service

	// database := client.Database(("palabras-express-api"))
	database := client.Database(("quickstart"))
	podcastsCollection := database.Collection("podcasts")

	service := handlers.Podcast{DB: podcastsCollection, Log: log}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(service.PodcastList),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
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
		return errors.Wrap(err, "listening and serving")

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "graceful shutdown")
		}
	}
	return nil
}

// func openDB() (*mongo.Client, error) {
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 	// formats the client
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

// 	return client, err
// }

// // structure to connect to the mongo db collections
// type Podcast struct {
// 	db *mongo.Collection
// }

// // PodcastList gets all the Podcast from teh db then encodes them in a response client
// func (p Podcast) PodcastList(w http.ResponseWriter, r *http.Request) {
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 	podcastList, err := podcast.List(p.db)
// 	if err != nil {
// 		panic(err)
// 	}

// 	podcastCursor, err := p.db.Find(ctx, bson.M{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err = podcastCursor.All(ctx, &podcastList); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(podcastList)

// 	data, err := json.Marshal(podcastList)
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
