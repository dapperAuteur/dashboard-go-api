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

// FilterTransactions gets all filtered transactions from the service layer.
// func (t Transaction) FilterTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	ctx, span := trace.StartSpan(ctx, "handlers.Transaction.FilterTransactions")
// 	defer span.End()

// 	var filterTranx budget.FilterTransaction

// 	if err := web.Decode(r, &filterTranx); err != nil {
// 		return err
// 	}

// 	list, err := budget.FilterTransactions(ctx, t.DB, filterTranx)
// 	if err != nil {
// 		return errors.Wrapf(err, "filtering transaction %q", filterTranx)
// 	}

// 	return web.Respond(ctx, w, list, http.StatusOK)
// }

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

// UpdateOneTransaction decodes the body of a request to update an existing transaction.
// The _id of the transaction is part of the request URL.
func (t *Transaction) UpdateOneTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	tranxID := chi.URLParam(r, "_id")

	var transactionUpdate budget.UpdateTransaction
	if err := web.Decode(r, &transactionUpdate); err != nil {
		return errors.Wrap(err, "decoding transaction update")
	}

	if err := budget.UpdateOneTransaction(ctx, t.DB, claims, tranxID, transactionUpdate, time.Now()); err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating transaction %q", tranxID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

// DeleteTransaction removes a single transaction identified by a transaction ID in the request URL.
func (t *Transaction) DeleteTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	tranxID := chi.URLParam(r, "_id")

	if err := budget.DeleteTransaction(ctx, t.DB, claims, tranxID); err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "deleting transaction %q", tranxID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
