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

// FinancialAccount defines all of the handlers related to FinancialAccount.
// It holds the application state needed by the handler methods.
type FinancialAccount struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// ListFinancialAccounts gets all the FinancialAccounts from the service layer.
func (fA FinancialAccount) ListFinancialAccounts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.FinancialAccount.ListFinancialAccounts")
	defer span.End()

	list, err := budget.ListFinancialAccounts(ctx, fA.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// RetrieveFinancialAccount will get the Financial Account from the db identified by an _id in the request URL, then encodes it in a response client.
func (fA FinancialAccount) RetrieveFinancialAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	faFound, err := budget.RetrieveFinancialAccount(ctx, fA.DB, _id)

	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for financial account %q", _id)
		}
	}

	return web.Respond(ctx, w, faFound, http.StatusOK)
}

// CreateFinancialAccount decodes the body of a reuqest to create a new financial account.
// The full budget with generated fields is sent back in the response.
func (fA FinancialAccount) CreateFinancialAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newFA budget.NewFinancialAccount

	if err := web.Decode(r, &newFA); err != nil {
		return err
	}

	faCreated, err := budget.CreateFinancialAccount(ctx, fA.DB, claims, newFA, time.Now())
	if err != nil {
		switch err {
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "creating financial account %q", newFA)
		}
	}

	return web.Respond(ctx, w, faCreated, http.StatusCreated)
}

// UpdateOneFinancialAccount decodes the body of a request to update an existing financial account.
// The _id of the financial account is part of the request URL.
func (fA *FinancialAccount) UpdateOneFinancialAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	finAccID := chi.URLParam(r, "_id")

	var finAccUpdate budget.UpdateFinancialAccount
	if err := web.Decode(r, &finAccUpdate); err != nil {
		return errors.Wrap(err, "decoding financial account update")
	}

	if err := budget.UpdateOneFinancialAccount(ctx, fA.DB, claims, finAccID, finAccUpdate, time.Now()); err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating financial account %q", finAccID)
		}
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}

// DeleteFinancialAccount removes a single financial account identified by a financial account ID in the request URL
func (fA *FinancialAccount) DeleteFinancialAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	finAccID := chi.URLParam(r, "_id")

	if err := budget.DeleteFinancialAccount(ctx, fA.DB, claims, finAccID); err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "deleting financial account %q", finAccID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
