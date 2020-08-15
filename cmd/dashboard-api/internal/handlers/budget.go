package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/budget"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Budget structure to connect to the mongo db collections
type Budget struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// List gets all the Budget from the service layer.
func (b Budget) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Budget.List")
	defer span.End()

	list, err := budget.List(ctx, b.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve get the Budget from the db identified by an _id in the request URL, then encodes it in a response client.
func (b Budget) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	budgetFound, err := budget.Retrieve(ctx, b.DB, _id)
	if err != nil {
		switch err {
		case budget.ErrBudgetNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case budget.ErrBudgetInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for budget %q", _id)
		}
	}
	return web.Respond(ctx, w, budgetFound, http.StatusOK)
}

// RetrieveByName gets the first Budget in the db with the provided name in the URL, then encodes it in a response client
func (b Budget) RetrieveByName(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	name := chi.URLParam(r, "name")

	budgetFound, err := budget.RetrieveByName(ctx, b.DB, name)
	if err != nil {
		switch err {
		case budget.ErrBudgetNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case budget.ErrBudgetInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for budget %q", name)
		}
	}
	return web.Respond(ctx, w, budgetFound, http.StatusOK)
}

// Create decodes the body of a request to create a new budget.
// The full budget with generated fields is sent back in the response.
func (b Budget) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newBudget budget.NewBudget

	if err := web.Decode(r, &newBudget); err != nil {
		return err
	}

	budget, err := budget.Create(ctx, b.DB, claims, newBudget, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, budget, http.StatusCreated)
}

// UpdateOne decodes the body of a request to update an existing budget.
// The _id of the budget is part of the request URL.
func (b *Budget) UpdateOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	budgetID := chi.URLParam(r, "_id")

	var budgetUpdate budget.UpdateBudget
	if err := web.Decode(r, &budgetUpdate); err != nil {
		return errors.Wrap(err, "decoding budget update")
	}

	if err := budget.UpdateOne(ctx, b.DB, claims, budgetID, budgetUpdate, time.Now()); err != nil {
		switch err {
		case budget.ErrBudgetNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case budget.ErrBudgetInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case budget.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating budget %q", budgetID)
		}
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}

// Delete removes a single budget identified by an budgetID in the request URL
func (b *Budget) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	budgetID := chi.URLParam(r, "_id")

	if err := budget.Delete(ctx, b.DB, claims, budgetID); err != nil {
		switch err {
		case budget.ErrBudgetNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case budget.ErrBudgetInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case budget.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "deleting budget %q", budgetID)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
