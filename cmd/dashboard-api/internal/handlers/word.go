package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/word"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

// Word defines all of the handlers related to words.
// It holds the application state needed by the handler methods.
type Word struct {
	DB  *mongo.Collection
	Log *log.Logger
}

// WordList gets all the Words from the service layer.
func (wd Word) WordList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Word.WordList")
	defer span.End()

	wordList, err := word.WordList(ctx, wd.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, wordList, http.StatusOK)
}
