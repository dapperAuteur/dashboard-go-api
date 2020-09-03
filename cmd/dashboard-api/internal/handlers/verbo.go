package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/word"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Verbo defines all of the handlers related to verbos.
// It holds the application state needed by the handler methods.
type Verbo struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// VerboList gets all the Verbos from the service layer.
func (v Verbo) VerboList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Verbo.VerboList")
	defer span.End()

	verboList, err := word.VerboList(ctx, v.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, verboList, http.StatusOK)
}

// CreateVerbo decodes the body of a request to create a new Verbo.
// The full Verbo with the generated fields is sent back in the response.
func (v Verbo) CreateVerbo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newVerbo word.NewVerbo

	if err := web.Decode(r, &newVerbo); err != nil {
		return err
	}

	verbo, err := word.CreateVerbo(ctx, v.DB, claims, newVerbo, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, verbo, http.StatusCreated)
}
