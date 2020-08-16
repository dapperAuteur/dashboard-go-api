package budget

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// List gets all the FinancialAccounts from the db then encodes them in a response client
func List(ctx context.Context, db *mongo.Collection) ([]FinancialAccount, error) {
	list := []FinancialAccount{}

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting cursor from financial accounts collection.")
	}

	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving financial accounts list")
	}

	return list, nil
}
