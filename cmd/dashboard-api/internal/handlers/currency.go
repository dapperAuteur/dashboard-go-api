package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/budget"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
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

// RetrieveCurrencyByID gets the Currency from the db identified by an _id in the request URL.
// Then encodes it in a response client.
func (c Currency) RetrieveCurrencyByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	currencyFound, err := budget.RetrieveCurrencyByID(ctx, c.DB, _id)
	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for currency %q", _id)
		}
	}

	return web.Respond(ctx, w, currencyFound, http.StatusOK)
}
