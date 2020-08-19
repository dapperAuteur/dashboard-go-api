package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/blog"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/pkg/errors"
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

// CreateNote decodes the body of a request to create a new note.
// The full note with generated fields is sent back in the response.
func (n Note) CreateNote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newNote blog.NewNote

	if err := web.Decode(r, &newNote); err != nil {
		return err
	}

	noteCreated, err := blog.CreateNote(ctx, n.DB, claims, newNote, time.Now())
	if err != nil {
		switch err {
		case apierror.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "creating vendor %q", newNote)
		}
	}

	return web.Respond(ctx, w, noteCreated, http.StatusCreated)
}
