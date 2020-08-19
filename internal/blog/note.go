package blog

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ListNotes gets all the Notes from the db then encodes them in a response client
func ListNotes(ctx context.Context, db *mongo.Collection) ([]Note, error) {

	list := []Note{}

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting cursor from note collection.")
	}

	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving note list")
	}
	return list, nil
}
