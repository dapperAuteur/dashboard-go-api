package word

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WordList gets all the Words form the db then encodes them in a response client.
func WordList(ctx context.Context, db *mongo.Collection) ([]Word, error) {

	wordList := []Word{}

	wordCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting wordCursor retrieving word list")
	}

	if err = wordCursor.All(ctx, &wordList); err != nil {
		return nil, errors.Wrapf(err, "retrieving word list")
	}

	return wordList, nil
}
