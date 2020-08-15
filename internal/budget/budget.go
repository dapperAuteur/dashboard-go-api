package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
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
