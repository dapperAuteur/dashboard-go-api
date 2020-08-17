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

// ListFinancialAccounts gets all the FinancialAccounts from the db then encodes them in a response client
func ListFinancialAccounts(ctx context.Context, db *mongo.Collection) ([]FinancialAccount, error) {
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

// CreateFinancialAccount takes data from the client to create a financial account in the db
func CreateFinancialAccount(ctx context.Context, db *mongo.Collection, user auth.Claims, newFA NewFinancialAccount, now time.Time) (*FinancialAccount, error) {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	financialAccount := FinancialAccount{
		AccountName:          newFA.AccountName,
		CurrentValue:         newFA.CurrentValue,
		FinancialInstitution: newFA.FinancialInstitution,
		MangerID:             user.Subject,
		CreatedAt:            now.UTC(),
		UpdatedAt:            now.UTC(),
	}

	faResult, err := db.InsertOne(ctx, financialAccount)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting financial account: %v", financialAccount)
	}

	fmt.Println("faResult : ", faResult)

	return &financialAccount, nil
}

// RetrieveFinancialAccount finds the financial account identified by a given _id.
func RetrieveFinancialAccount(ctx context.Context, db *mongo.Collection, _id string) (*FinancialAccount, error) {

	var financialAccount FinancialAccount

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&financialAccount); err != nil {
		return nil, apierror.ErrNotFound
	}

	return &financialAccount, nil
}

// UpdateOneFinancialAccount modifies data about a Financial Account.
// It will error if the specified _id is invalid or does NOT reference an existing Financial Account.
func UpdateOneFinancialAccount(ctx context.Context, db *mongo.Collection, user auth.Claims, faID string, updateFA UpdateFinancialAccount, now time.Time) error {

	faObjectID, err := primitive.ObjectIDFromHex(faID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundFA, err := RetrieveFinancialAccount(ctx, db, faID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("financial account to update found %+v : \n", foundFA)

	var (
		isAdmin = user.HasRole(auth.RoleAdmin)
		isOwner = foundFA.MangerID == user.Subject
		canView = isAdmin || isOwner
	)

	if !canView {
		return apierror.ErrForbidden
	}

	financialAccount := FinancialAccount{}

	if updateFA.AccountName != nil {
		financialAccount.AccountName = *updateFA.AccountName
	}

	if updateFA.CurrentValue != nil {
		financialAccount.CurrentValue = *updateFA.CurrentValue
	}

	if updateFA.FinancialInstitution != nil {
		financialAccount.FinancialInstitution = *updateFA.FinancialInstitution
	}

	financialAccount.ID = faObjectID

	financialAccount.UpdatedAt = now

	updateFinAcc := bson.M{
		"$set": financialAccount,
	}

	faResult, err := db.UpdateOne(ctx, bson.M{"_id": faObjectID}, updateFinAcc)
	if err != nil {
		return errors.Wrap(err, "updating financial account")
	}

	fmt.Printf("faResult updated %v : \n", faResult)

	return nil
}

// DeleteFinancialAccount removes the financial account identified by a given _id
func DeleteFinancialAccount(ctx context.Context, db *mongo.Collection, user auth.Claims, faID string) error {

	faObjectID, err := primitive.ObjectIDFromHex(faID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundFA, err := RetrieveFinancialAccount(ctx, db, faID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("financial account to delete found %+v : \n", foundFA)

	var (
		isAdmin = user.HasRole(auth.RoleAdmin)
		isOwner = foundFA.MangerID == user.Subject
		canView = isAdmin || isOwner
	)

	if !canView {
		return apierror.ErrForbidden
	}

	result, err := db.DeleteOne(ctx, bson.M{"_id": faObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting financial account %s", faID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
