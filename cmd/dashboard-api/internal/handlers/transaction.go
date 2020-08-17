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
