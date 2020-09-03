package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/word"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Affix defines all of the handlers related to affixes.
// It holds the application state needed by the handler methods.
type Affix struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// AffixList gets all the Affixes from the service layer.
func (a Affix) AffixList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Affix.AffixList")
	defer span.End()

	affixList, err := word.AffixList(ctx, a.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, affixList, http.StatusOK)
}

// RetrieveAffixByID gets the Affix from the db identified by an _id in the request URL, then encodes it in a response client.
func (a Affix) RetrieveAffixByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	affixFound, err := word.RetrieveAffixByID(ctx, a.DB, _id)
	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for affix %q", _id)
		}
	}

	return web.Respond(ctx, w, affixFound, http.StatusOK)
}
