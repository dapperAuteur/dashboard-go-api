package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
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

// CreateCurrency adds a Currency to the database.
// It returns the created Currency with the fields populated.
func CreateCurrency(ctx context.Context, db *mongo.Collection, user auth.Claims, newCurrency NewCurrency, now time.Time) (*Currency, error) {

	isAdmin := user.HasRole(auth.RoleAdmin)
	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	currency := Currency{
		CurrencyName: newCurrency.CurrencyName,
		CurrencyType: newCurrency.CurrencyType,
		Symbol:       newCurrency.Symbol,
		CreatedAt:    now.UTC(),
		UpdatedAt:    now.UTC(),
	}

	currencyResult, err := db.InsertOne(ctx, currency)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Currency: %v", newCurrency)
	}

	fmt.Println("currencyResult : ", currencyResult)

	return &currency, nil
}

// UpdateOneCurrency modifies data about one Currency.
// It will ERROR if the specified currencyID is invalid or does NOT reference existing currency.
func UpdateOneCurrency(ctx context.Context, db *mongo.Collection, user auth.Claims, currencyID string, updateCurrency UpdateCurrency, now time.Time) error {

	currencyObjectID, err := primitive.ObjectIDFromHex(currencyID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundCurrency, err := RetrieveCurrencyByID(ctx, db, currencyID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("currency to update found %+v : \n", foundCurrency)

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	currency := Currency{}

	if updateCurrency.CurrencyName != nil {
		currency.CurrencyName = *updateCurrency.CurrencyName
	}

	if updateCurrency.CurrencyType != nil {
		currency.CurrencyType = *updateCurrency.CurrencyType
	}

	if updateCurrency.Symbol != nil {
		currency.Symbol = *updateCurrency.Symbol
	}

	currency.ID = currencyObjectID

	currency.UpdatedAt = now

	updateC := bson.M{
		"$set": currency,
	}

	currencyResult, err := db.UpdateOne(ctx, bson.M{"_id": currencyObjectID}, updateC)
	if err != nil {
		return errors.Wrap(err, "updating currency")
	}

	fmt.Printf("currencyResult updated %v : \n", currencyResult)

	return nil
}

// DeleteCurrencyByID removes the Currency identified by a given ID.
func DeleteCurrencyByID(ctx context.Context, db *mongo.Collection, user auth.Claims, currencyID string) error {

	currencyObjectID, err := primitive.ObjectIDFromHex(currencyID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	_, err = RetrieveCurrencyByID(ctx, db, currencyID)
	if err != nil {
		return apierror.ErrNotFound
	}

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	result, err := db.DeleteOne(ctx, bson.M{"_id": currencyObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting currency %s", currencyID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
