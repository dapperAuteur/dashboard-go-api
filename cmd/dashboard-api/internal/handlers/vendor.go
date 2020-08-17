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

// Vendor defines all of the handlers related to vendor.
// It holds the application state needed by the handler methods.
type Vendor struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// ListVendors gets all vendors from the service layer.
func (v Vendor) ListVendors(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Vendor.ListVendors")
	defer span.End()

	list, err := budget.ListVendors(ctx, v.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}
