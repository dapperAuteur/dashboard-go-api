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

// UpdateOneTransaction modifies data about a transaction.
// It will error if the specified _id is invalid or does NOT reference an existing transaction.
func UpdateOneTransaction(ctx context.Context, db *mongo.Collection, user auth.Claims, tranxID string, updateTranx UpdateTransaction, now time.Time) error {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	tObjectID, err := primitive.ObjectIDFromHex(tranxID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundTranx, err := RetrieveTransaction(ctx, db, tranxID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("transaction to update found %+v : \n", foundTranx)

	transaction := Transaction{}

	if updateTranx.BudgetID != nil {
		transaction.BudgetID = *updateTranx.BudgetID
	}

	if updateTranx.CurrencyID != nil {
		transaction.CurrencyID = *updateTranx.CurrencyID
	}

	if updateTranx.FinancialAccountID != nil {
		transaction.FinancialAccountID = *updateTranx.FinancialAccountID
	}

	// if updateTranx.Occurrence != nil {
	// 	transaction.Occurrence = *updateTranx.Occurrence
	// }

	if updateTranx.TransactionEvent != nil {
		transaction.TransactionEvent = *updateTranx.TransactionEvent
	}

	if updateTranx.TransactionValue != nil {
		transaction.TransactionValue = *updateTranx.TransactionValue
	}

	if updateTranx.VendorID != nil {
		transaction.VendorID = *updateTranx.VendorID
	}

	if updateTranx.ParticipantID != nil {
		transaction.ParticipantID = *updateTranx.ParticipantID
	}

	transaction.ID = tObjectID

	transaction.UpdatedAt = now

	updateTransaction := bson.M{
		"$set": transaction,
	}

	tranxResult, err := db.UpdateOne(ctx, bson.M{"_id": tObjectID}, updateTransaction)
	if err != nil {
		return errors.Wrap(err, "updating transaction")
	}

	fmt.Printf("tranxResult updated %v : \n", tranxResult)

	return nil
}
