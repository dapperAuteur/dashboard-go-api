package main

import (
	"context"
	"crypto/rsa"
	_ "expvar" // Register the expvar handlers
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof" // register the /debug/pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/dapperAuteur/dashboard-go-api/cmd/dashboard-api/internal/handlers"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/conf"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/database"
	jwt "github.com/dgrijalva/jwt-go"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
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

	dbCredentials := os.Getenv("DATABASE_CRED")
	fmt.Printf("************\n dbCredentials", dbCredentials)
	a := `conf:"default:mongodb+srv://` + dbCredentials + `@localhost/palabras-express-api?retryWrites=true&w=majority"`
	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8080"`
			Debug           string        `conf:"default:localhost:6060"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			AtlasURI string `a` // connection string for docker Mongo image
			// AtlasURI string `conf:"default:mongodb+srv://awe:XjtsRQPAjyDbokQE@localhost/palabras-express-api?retryWrites=true&w=majority"` // connection string for docker Mongo image
		}
		Auth struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		Trace struct {
			URL         string  `conf:"default:http://localhost:9411/api/v2/spans"`
			Service     string  `conf:"default:dashboard-api"`
			Probability float64 `conf:"default:1"`
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

	// is it ok to do this twice, I think ctx is needed here to close the context later
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// ==
	// Initialize authentication support
	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}
	// =
	// Start Database

	client, err := database.Open(database.Config{
		AtlasURI: cfg.DB.AtlasURI,
	})
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// ==
	// Start Tracing Support

	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer closer()

	// ==
	// Start Debug Service
	go func() {
		log.Printf("main : Debug service listening on %s", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Printf("main : Debug service ended %v", err)
	}()

	// =========================================================================
	// Start API Service

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// send the db to the handler and let the router determine which collection to use
	myDatabase := client.Database(("quickstart")) // development database
	// myDatabase := client.Database(("palabras-express-api")) // production database

	// service := handlers.Podcast{DB: podcastsCollection, Log: log}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(shutdown, log, myDatabase, authenticator),
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

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case sig := <-shutdown:
		log.Println("main : Start shutdown", sig)

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

		if sig == syscall.SIGSTOP {
			return errors.New("integrity error detectecd, asking for self shutdown")
		}
	}
	return nil
}

func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}

func registerTracer(service, httpAddr, traceURL string, probabilty float64) (func() error, error) {
	localEndpoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probabilty),
	})
	return reporter.Close, nil

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
