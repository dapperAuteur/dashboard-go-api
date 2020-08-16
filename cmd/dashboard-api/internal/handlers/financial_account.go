package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// FinancialAccount structure to connect to the mongo db collection
type FinancialAccount struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// ListFinancialAccounts gets all the FinancialAccounts from the service layer.
func (fA FinancialAccount) ListFinancialAccounts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.FinancialAccount.List")
	defer span.End()

	list, err := FinancialAccount.List(ctx, fA.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}
