package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/word"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
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

// RetrieveWordByID gets the Word from the db identified by an _id in the request URL, then encodes it in a response client.
func (wd Word) RetrieveWordByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_id := chi.URLParam(r, "_id")

	wordFound, err := word.RetrieveWordByID(ctx, wd.DB, _id)
	if err != nil {
		switch err {
		case apierror.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case apierror.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for word %q", _id)
		}
	}

	return web.Respond(ctx, w, wordFound, http.StatusOK)
}

// RetrieveWord gets the Word from the db identified by an word string in the request URL, then encodes it in a response client.
// func (wd Word) RetrieveWord(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	word := chi.URLParam(r, "word")

// 	wordFound, err := word.RetrieveWord(ctx, wd.DB, word)
// 	if err != nil {
// 		switch err {
// 		case apierror.ErrNotFound:
// 			return web.NewRequestError(err, http.StatusNotFound)
// 		case apierror.ErrInvalidID:
// 			return web.NewRequestError(err, http.StatusBadRequest)
// 		default:
// 			return errors.Wrapf(err, "looking for word %q", word)
// 		}
// 	}

// 	return web.Respond(ctx, w, wordFound, http.StatusOK)
// }

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
