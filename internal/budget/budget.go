package budget

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/mongo"
)

// List gets all the Budgets from the db then encodes them in a response client
func List(ctx context.Context, db *mongo.Collection) ([]Budget, error) {
	list := []Budget{}

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting cursor from budget db. retrieving budget list")
	}

	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving budget list")
	}

	return list, nil

}
