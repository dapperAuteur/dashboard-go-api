package budget

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ListTransactions gets all the Transactions from the db then encodes them in a response client.
func ListTransactions(ctx context.Context, db *mongo.Collection) ([]Transaction, error) {

	list := []Transaction{}

	cursor, err := db.Find(ctx, bson.M{})
	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving transaction list")
	}

	return list, nil
}
