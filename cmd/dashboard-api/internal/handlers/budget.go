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
