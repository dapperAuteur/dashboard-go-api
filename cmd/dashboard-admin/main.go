package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/conf"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/database"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// ==
	// Configuration

	var cfg struct {
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

	// print config values when app starts
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// ==
	// Start Database
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := database.Open(database.Config{
		AtlasUri: cfg.DB.AtlasUri,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer client.Disconnect(ctx)

	return nil

}

// Do we have a use for a switch statement to run admin commands?

// func openDB() (*mongo.Client, error) {
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 	// formats the client
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://awe:XjtsRQPAjyDbokQE@palabras-express-api.whbeh.mongodb.net/palabras-express-api?retryWrites=true&w=majority"))

// 	return client, err
// }
