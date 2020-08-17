package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Retrieve finds the budget identified by a given _id.
func Retrieve(ctx context.Context, db *mongo.Collection, _id string) (*Budget, error) {

	var budget Budget

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&budget); err != nil {
		return nil, apierror.ErrNotFound
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

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	budget := Budget{
		ManagerID:   user.Subject,
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

// UpdateOne modifies data about a Budget.
// It will error if the specified _id is invalid or does NOT reference an existing Budget.
func UpdateOne(ctx context.Context, db *mongo.Collection, user auth.Claims, budgetID string, updateBudget UpdateBudget, now time.Time) error {

	budgetObjectID, err := primitive.ObjectIDFromHex(budgetID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundBudget, err := Retrieve(ctx, db, budgetID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("budget to update found %+v : \n", foundBudget)

	var (
		isAdmin = user.HasRole(auth.RoleAdmin)
		isOwner = foundBudget.ManagerID == user.Subject
		canView = isAdmin || isOwner
	)

	if !canView {
		return apierror.ErrForbidden
	}

	budget := Budget{}

	if updateBudget.BudgetName != nil {
		budget.BudgetName = *updateBudget.BudgetName
	}

	if updateBudget.BudgetValue != nil {
		budget.BudgetValue = *updateBudget.BudgetValue
	}

	budget.ID = budgetObjectID

	budget.UpdatedAt = now

	updateB := bson.M{
		"$set": budget,
	}

	budgetResult, err := db.UpdateOne(ctx, bson.M{"_id": budgetObjectID}, updateB)
	if err != nil {
		return errors.Wrap(err, "updating budget")
	}

	fmt.Printf("budgetResult updated %v : \n", budgetResult)

	return nil
}

// Delete removes the budget identified by a given _id
func Delete(ctx context.Context, db *mongo.Collection, user auth.Claims, budgetID string) error {

	budgetObjectID, err := primitive.ObjectIDFromHex(budgetID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundBudget, err := Retrieve(ctx, db, budgetID)
	if err != nil {
		return apierror.ErrNotFound
	}

	var (
		isAdmin = user.HasRole(auth.RoleAdmin)
		isOwner = foundBudget.ManagerID == user.Subject
		canView = isAdmin || isOwner
	)

	if !canView {
		return apierror.ErrForbidden
	}

	result, err := db.DeleteOne(ctx, bson.M{"_id": budgetObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting budget %s", budgetID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
