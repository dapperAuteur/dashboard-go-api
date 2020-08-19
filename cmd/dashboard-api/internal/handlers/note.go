package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/blog"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Note defines all of the handlers related to note.
// It holds the application state needed by the handler methods.
type Note struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// ListNotes gets all notes from the service layer.
func (n Note) ListNotes(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Note.ListNotes")
	defer span.End()

	list, err := blog.ListNotes(ctx, n.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}
