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

	// Finance Related
	budgetsCollection := db.Collection("budgets")
	financialAccountsCollection := db.Collection("financialaccounts")
	vendorsCollection := db.Collection("vendors")
	transactionsCollection := db.Collection("transactions")

	// Podcast Related
	episodesCollection := db.Collection("episodes")
	podcastsCollection := db.Collection("podcasts")

	// Finance Related
	budget := Budget{
		DB:  budgetsCollection,
		Log: logger,
	}

	financialAccount := FinancialAccount{
		DB:  financialAccountsCollection,
		Log: logger,
	}

	transaction := Transaction{
		DB:  transactionsCollection,
		Log: logger,
	}

	vendor := Vendor{
		DB:  vendorsCollection,
		Log: logger,
	}

	// Content Creation

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
	app.Handle(http.MethodGet, "/v1/budgets/{_id}", budget.Retrieve)
	app.Handle(http.MethodGet, "/v1/budgets/{name}", budget.RetrieveByName)
	app.Handle(http.MethodPost, "/v1/budgets", budget.Create, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/v1/budgets/{_id}", budget.UpdateOne, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/v1/budgets/{_id}", budget.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// FinancialAccount Routes
	app.Handle(http.MethodGet, "/v1/financial-accounts", financialAccount.ListFinancialAccounts)
	app.Handle(http.MethodPost, "/v1/financial-accounts", financialAccount.CreateFinancialAccount, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/financial-accounts/{_id}", financialAccount.RetrieveFinancialAccount)
	app.Handle(http.MethodPut, "/v1/financial-accounts/{_id}", financialAccount.UpdateOneFinancialAccount, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/v1/financial-accounts/{_id}", financialAccount.DeleteFinancialAccount, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// Transaction Routes
	app.Handle(http.MethodGet, "/v1/transactions", transaction.ListTransactions)
	app.Handle(http.MethodPost, "/v1/transactions", transaction.CreateTransaction, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/transactions/{_id}", transaction.RetrieveTransaction)

	// Vendor Routes
	app.Handle(http.MethodGet, "/v1/vendors", vendor.ListVendors)
	app.Handle(http.MethodPost, "/v1/vendors", vendor.CreateVendor, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/vendors/{_id}", vendor.RetrieveVendor)
	app.Handle(http.MethodPut, "/v1/vendors/{_id}", vendor.UpdateOneVendor, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/v1/vendors/{_id}", vendor.DeleteVendor, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// Episode Routes
	app.Handle(http.MethodGet, "/v1/episodes", episode.EpisodeList)
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}/episodes", episode.PodcastEpisodeList)
	app.Handle(http.MethodGet, "/v1/episodes/{episodeID}", episode.RetrieveEpisode)
	app.Handle(http.MethodPost, "/v1/podcasts/{_id}/episodes", episode.AddEpisode, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	// app.Handle(http.MethodGet, "/v1/podcasts/{_id}/episodes/{_id}", episode.Retrieve, mid.Authenticate(authenticator))

	// Podcast Routes
	app.Handle(http.MethodGet, "/v1/podcasts", podcast.PodcastList)
	app.Handle(http.MethodPost, "/v1/podcasts", podcast.CreatePodcast, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/podcasts/{_id}", podcast.Retrieve)
	app.Handle(http.MethodGet, "/v1/podcasts/{title}", podcast.Retrieve)
	app.Handle(http.MethodPut, "/v1/podcasts/{_id}", podcast.UpdateOnePodcast, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/v1/podcasts/{_id}", podcast.DeletePodcast, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	return app
}
