package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/budget"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
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
