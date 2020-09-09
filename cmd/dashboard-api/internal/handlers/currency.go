package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/budget"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Currency defines all of the handlers related to currencies.
// It holds the application state needed by the handler methods.
type Currency struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// CurrencyList gets all the Currencies from the service layer.
func (c Currency) CurrencyList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Currency.CurrencyList")
	defer span.End()

	currencyList, err := budget.CurrencyList(ctx, c.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, currencyList, http.StatusOK)
}
