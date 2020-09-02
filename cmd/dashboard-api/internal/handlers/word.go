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

// RetrieveWord gets the Word from the db identified by an _id in the request URL, then encodes it in a response client.

// CreateWord decodes the body of a request to create a new Word.
// The full Word with generated fields is sent back in the response.
func (wd Word) CreateWord(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims missing from context")
	}

	var newWord word.NewWord

	if err := web.Decode(r, &newWord); err != nil {
		return err
	}

	word, err := word.CreateWord(ctx, wd.DB, claims, newWord, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, word, http.StatusCreated)
}
