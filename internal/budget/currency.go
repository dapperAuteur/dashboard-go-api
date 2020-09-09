package budget

import (
	"context"
	"fmt"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CurrencyList gets all the currencies from the database then encodes them in a response client.
func CurrencyList(ctx context.Context, db *mongo.Collection) ([]Currency, error) {

	currencyList := []Currency{}

	currencyCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting currencyCursor retrieving currency list")
	}

	if err = currencyCursor.All(ctx, &currencyList); err != nil {
		return nil, errors.Wrapf(err, "retrieving currency list")
	}

	return currencyList, nil
}

// RetrieveCurrencyByID gets the first Currency in the db with the provided ID.
func RetrieveCurrencyByID(ctx context.Context, db *mongo.Collection, _id string) (*Currency, error) {

	var currency Currency

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&currency); err != nil {
		return nil, apierror.ErrNotFound
	}

	fmt.Println("currency found : ", currency)

	return &currency, nil
}
