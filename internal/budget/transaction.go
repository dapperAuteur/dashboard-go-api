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

// ListTransactions gets all the Transactions from the db then encodes them in a response client.
func ListTransactions(ctx context.Context, db *mongo.Collection) ([]Transaction, error) {

	list := []Transaction{}

	cursor, err := db.Find(ctx, bson.M{})
	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving transaction list")
	}

	return list, nil
}

// CreateTransaction takes data from the client to create a transaction in the db
func CreateTransaction(ctx context.Context, db *mongo.Collection, user auth.Claims, newTranx NewTransaction, now time.Time) (*Transaction, error) {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	tranx := Transaction{
		BudgetID:           newTranx.BudgetID,
		CurrencyID:         newTranx.CurrencyID,
		FinancialAccountID: newTranx.FinancialAccountID,
		// Occurrence:         now.UTC(),
		TransactionEvent: newTranx.TransactionEvent,
		TransactionValue: newTranx.TransactionValue,
		VendorID:         newTranx.VendorID,
		ParticipantID:    newTranx.ParticipantID,
		CreatedAt:        now.UTC(),
		UpdatedAt:        now.UTC(),
	}

	tranxResult, err := db.InsertOne(ctx, tranx)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting transaction : %v", tranx)
	}

	fmt.Println("tranxResult : ", tranxResult)

	return &tranx, nil
}

// RetrieveTransaction finds a single Transaction by _id
func RetrieveTransaction(ctx context.Context, db *mongo.Collection, _id string) (*Transaction, error) {

	var transaction Transaction

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction); err != nil {
		return nil, apierror.ErrNotFound
	}

	return &transaction, nil
}
