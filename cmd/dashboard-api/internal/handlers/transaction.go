package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/budget"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Transaction defines all of the handlers related to transaction.
// It holds the application state needed by the handler methods.
type Transaction struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// ListTransactions gets all transactions from the service layer.
func (t Transaction) ListTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Transaction.ListTransactions")
	defer span.End()

	list, err := budget.ListTransactions(ctx, t.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// CreateTransaction decodes the body of a request to create a new transaction.
// The full transaction with generated fields is sent back in the response.
func (t Transaction) CreateTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newTransaction budget.NewTransaction

	if err := web.Decode(r, &newTransaction); err != nil {
		return err
	}

	tranxCreated, err := budget.CreateTransaction(ctx, t.DB, claims, newTransaction, time.Now())
	if err != nil {
		switch err {
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "creating transaction %q", newTransaction)
		}
	}

	return web.Respond(ctx, w, tranxCreated, http.StatusCreated)
}

// RetrieveTransaction will get the tranx from the db identified by an _id in the request URL, then encodes it in a response client.
func (t Transaction) RetrieveTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	tranxFound, err := budget.RetrieveTransaction(ctx, t.DB, _id)
	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for transaction %q", _id)
		}
	}

	return web.Respond(ctx, w, tranxFound, http.StatusOK)
}
