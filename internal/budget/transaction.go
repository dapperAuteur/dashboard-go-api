package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/utility"
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

// FilterTransactions gets all the Transactions from the database that filter the filter criteria.
// The encodes them in a response client. It returns an empty array if none fit the criteria.
// func FilterTransactions(ctx context.Context, db *mongo.Collection, filterTranx FilterTransaction) ([]Transaction, error) {

// 	list := []Transaction{}
// 	filter := bson.M{}

// 	if filterTranx.BudgetID != nil {
// 		filter.BudgetID = filterTranx.BudgetID
// 	}

// 	if filterTranx.CurrencyID != nil {
// 		filter.CurrencyID = filterTranx.CurrencyID
// 	}

// 	if filterTranx.FinancialAccountID != nil {
// 		filter.FinancialAccountID = filterTranx.FinancialAccountID
// 	}

// 	if filterTranx.OccurrenceString != nil {
// 		filter.OccurrenceString = filterTranx.OccurrenceString
// 	}

// 	if filterTranx.TransactionEvent != nil {
// 		filter.TransactionEvent = filterTranx.TransactionEvent
// 	}

// 	if filterTranx.TransactionCredit != nil {
// 		filter.TransactionCredit = filterTranx.TransactionCredit
// 	}

// 	if filterTranx.TransactionDebit != nil {
// 		filter.TransactionDebit = filterTranx.TransactionDebit
// 	}

// 	if filterTranx.VendorID != nil {
// 		filter.VendorID = filterTranx.VendorID
// 	}

// 	if filterTranx.ParticipantID != nil {
// 		filter.ParticipantID = filterTranx.ParticipantID
// 	}

// 	if filterTranx.CreatedAt != nil {
// 		filter.CreatedAt = filterTranx.CreatedAt
// 	}

// 	if filterTranx.UpdatedAt != nil {
// 		filter.UpdatedAt = filterTranx.UpdatedAt
// 	}

// 	cursor, err := db.Find(ctx, filter)
// 	if err = cursor.All(ctx, &list); err != nil {
// 		return nil, errors.Wrapf(err, "filtering transactions")
// 	}

// 	return list, nil
// }

// CreateTransaction takes data from the client to create a transaction in the db
func CreateTransaction(ctx context.Context, db *mongo.Collection, user auth.Claims, newTranx NewTransaction, now time.Time) (*Transaction, error) {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	tranx := Transaction{}

	var (
		// finAcctObjectIDs, participantObjectIDs []primitive.ObjectID
		finAcctIDsSlice, participantIDsSlice []string
	)

	// check if prop is provided
	if newTranx.FinancialAccountID != nil {
		// convert []newTranx.FinancialAccountID (ObjectID) to []string
		// objIDs, err := utility.SliceStringsToObjectIDs(*newTranx.FinancialAccountID)
		// if err != nil {
		// 	return nil, err
		// }
		objIDs := append(*newTranx.FinancialAccountID, *newTranx.FinancialAccountID...)
		finAcctIDsSlice = utility.RemoveDuplicateStringValues(objIDs)
	}

	// check if prop is provided
	if newTranx.ParticipantID != nil {
		// convert []newTranx.ParticipantID (ObjectID) to []string
		// objIDs, err := utility.SliceStringsToObjectIDs(*newTranx.ParticipantID)
		// if err != nil {
		// 	return nil, err
		// }
		objIDs := append(*newTranx.ParticipantID, *newTranx.ParticipantID...)
		// participantObjectIDs = utility.RemoveDuplicateObjectIDValues(objIDs)
		participantIDsSlice = utility.RemoveDuplicateStringValues(objIDs)
	}

	tranx = Transaction{
		BudgetID:           newTranx.BudgetID,
		CurrencyID:         newTranx.CurrencyID,
		FinancialAccountID: finAcctIDsSlice,
		// Occurrence:         now.UTC(),
		OccurrenceString:  newTranx.OccurrenceString,
		TransactionEvent:  newTranx.TransactionEvent,
		TransactionCredit: newTranx.TransactionCredit,
		TransactionDebit:  newTranx.TransactionDebit,
		VendorID:          newTranx.VendorID,
		ParticipantID:     participantIDsSlice,
		CreatedAt:         now.UTC(),
		UpdatedAt:         now.UTC(),
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

	// fmt.Println("&transaction.FinancialAccountID", &transaction.FinancialAccountID)
	// fmt.Printf("***************\n&transaction.FinancialAccountID Type : %T\n", &transaction.FinancialAccountID)

	return &transaction, nil
}

// UpdateOneTransaction modifies data about a transaction.
// It will error if the specified _id is invalid or does NOT reference an existing transaction.
func UpdateOneTransaction(ctx context.Context, db *mongo.Collection, user auth.Claims, tranxID string, updateTranx UpdateTransaction, now time.Time) error {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	foundTranx, err := RetrieveTransaction(ctx, db, tranxID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("transaction to update found %+v : \n", foundTranx)

	tObjectID, err := primitive.ObjectIDFromHex(tranxID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	transaction := Transaction{}

	if updateTranx.BudgetID != nil {
		transaction.BudgetID = *updateTranx.BudgetID
	}

	if updateTranx.CurrencyID != nil {
		transaction.CurrencyID = *updateTranx.CurrencyID
	}

	if updateTranx.FinancialAccountID != nil {
		// take *updateTranx.FinancialAccountID.
		// convert to []primitive.ObjectID and return
		// finAcctObjectIDs, err := utility.SliceStringsToObjectIDs(*updateTranx.FinancialAccountID)
		objectIDs := append(*updateTranx.FinancialAccountID, foundTranx.FinancialAccountID...)
		uniqueFinAccObjIDs := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}

		transaction.FinancialAccountID = uniqueFinAccObjIDs
	}

	// if updateTranx.Occurrence != nil {
	// 	transaction.Occurrence = *updateTranx.Occurrence
	// }

	if updateTranx.OccurrenceString != nil {
		transaction.OccurrenceString = *updateTranx.OccurrenceString
	}

	if updateTranx.TransactionEvent != nil {
		transaction.TransactionEvent = *updateTranx.TransactionEvent
	}

	if updateTranx.TransactionCredit != nil {
		transaction.TransactionCredit = *updateTranx.TransactionCredit
	}

	if updateTranx.TransactionDebit != nil {
		transaction.TransactionDebit = *updateTranx.TransactionDebit
	}

	if updateTranx.VendorID != nil {
		transaction.VendorID = *updateTranx.VendorID
	}

	if updateTranx.ParticipantID != nil {
		// participantObjectIDs, err := utility.SliceStringsToObjectIDs(*updateTranx.ParticipantID)
		objectIDs := append(*updateTranx.ParticipantID, foundTranx.ParticipantID...)
		uniquePartObjIDs := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}
		transaction.ParticipantID = uniquePartObjIDs
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

// DeleteTransaction removes the transaction identified by a given _id
func DeleteTransaction(ctx context.Context, db *mongo.Collection, user auth.Claims, tranxID string) error {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	tranxObjectID, err := primitive.ObjectIDFromHex(tranxID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundTranx, err := RetrieveTransaction(ctx, db, tranxID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("transaction to delelete found %+v : \n", foundTranx)

	result, err := db.DeleteOne(ctx, bson.M{"_id": tranxObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting transaction %s", tranxID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
