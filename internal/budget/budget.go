package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Predefined Errors indentify expected failure conditions.
var (
	// ErrBudgetNotFound is used when a specific Budget is requested but does not exist.
	ErrBudgetNotFound = errors.New("budget NOT found")

	// ErrBudgetInvalID is used when an invalid ID is provided.
	ErrBudgetInvalID = errors.New("_id is NOT in its proper form")
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

// Get finds the budget identified by a given _id.
func Retrieve(ctx context.Context, db *mongo.Collection, _id string) (*Budget, error) {

	var budget Budget

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, ErrBudgetInvalID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&budget); err != nil {
		return nil, ErrBudgetNotFound
	}

	return &budget, nil
}

// RetrieveByName gets the first Budget in the db identified by a given name
func RetrieveByName(ctx context.Context, db *mongo.Collection, name string) (*Budget, error) {

	var budget Budget

	filter := Budget{
		BudgetName: name,
	}

	if err := db.FindOne(ctx, filter).Decode(&budget); err != nil {
		return nil, errors.Wrapf(err, "retrieving budget by name: %s", name)
	}

	return &budget, nil
}

// Create adds a Budget to the database.
// It returns the created Budget with fields like ID and CreatedAt populated.
func Create(ctx context.Context, db *mongo.Collection, user auth.Claims, newBudget NewBudget, now time.Time) (*Budget, error) {

	budget := Budget{
		BudgetName:  newBudget.BudgetName,
		BudgetValue: newBudget.BudgetValue,
		CreatedAt:   now.UTC(),
		UpdatedAt:   now.UTC(),
	}

	budgetResult, err := db.InsertOne(ctx, budget)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Budget: %v", newBudget)
	}
	fmt.Println("podcastResult : ", budgetResult)

	return &budget, nil
}
