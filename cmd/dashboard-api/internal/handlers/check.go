package handlers

import (
	"context"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/database"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
)

// Check has handlers to implement service orchestration
// add logger later
type Check struct {
	DB *mongo.Collection
}

// Health responds with a 200 OK if the service is healthy and ready for traffic.
func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}
	if err := database.StatusCheck(ctx, c.DB); err != nil {
		health.Status = "db NOT ready"
		return web.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "OK"
	return web.Respond(ctx, w, health, http.StatusOK)
}
