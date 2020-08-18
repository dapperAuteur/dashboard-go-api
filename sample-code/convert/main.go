package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	// create a nil slice of string and ObjectID
	// var sliceObjectIDs []primitive.ObjectID
	var sliceStrings []string

	// add strings to string slice
	sliceStrings = append(sliceStrings, "5f39f4465f148fdb7c5d6e9f", "5f39f45a5f148fdb7c5d6ea0", "5f39f4635f148fdb7c5d6ea1", "5f39f46c5f148fdb7c5d6ea2", "5f39f4745f148fdb7c5d6ea3", "5f39f47e5f148fdb7c5d6ea4", "5f39f48c5f148fdb7c5d6ea5", "5f39f4965f148fdb7c5d6ea6", "5f39f4a15f148fdb7c5d6ea7", "5f39f4ab5f148fdb7c5d6ea8", "5f39f4b55f148fdb7c5d6ea9", "5f39f4bd5f148fdb7c5d6eaa", "5f39f4c65f148fdb7c5d6eab", "5f39f4d55f148fdb7c5d6eac", "5f39f5365f148fdb7c5d6ead")

	oIDs, err := SliceStringsToObjectIDs(sliceStrings)
	if err != nil {
		return
	}
	// sliceObjectIDs = append(sliceObjectIDs)

	fmt.Println("oIDs", oIDs)
	// fmt.Println("sliceStrings", sliceStrings)
	// // check the ids of params are provided by client
	// if sliceStrings != nil {
	// 	// loop thru *newTranx.FinancialAccountID and append them to finAcctObjectIDs
	// 	for _, stringID := range sliceStrings {
	// 		// confirm each string id can be found in the db
	// 		// _, err := RetrieveFinancialAccount(ctx, db, stringID)
	// 		// if err != nil {
	// 		// 	return apierror.ErrNotFound
	// 		// }

	// 		fmt.Println("stringID", stringID)
	// 		// convert/check each element to confirm its a primitive.ObjectID before adding to the slice
	// 		objectID, err := primitive.ObjectIDFromHex(stringID)
	// 		if err != nil {
	// 			return
	// 		}
	// 		sliceObjectIDs = append(sliceObjectIDs, objectID)

	// 	}

	// }
	// fmt.Println("sliceObjectIDs", sliceObjectIDs)

}

// SliceStringsToObjectIDs takes a slice of strings and converts them to a slice of ObjectIDs
func SliceStringsToObjectIDs(sliceStrings []string) ([]primitive.ObjectID, error) {

	var sliceObjectIDs []primitive.ObjectID

	if sliceStrings != nil {

		for _, stringID := range sliceStrings {

			fmt.Println("stringID : %s\n", stringID)
			objectID, err := primitive.ObjectIDFromHex(stringID)
			if err != nil {
				return nil, err
			}
			sliceObjectIDs = append(sliceObjectIDs, objectID)
		}
	}
	fmt.Println("sliceObjectIDs", sliceObjectIDs)

	return sliceObjectIDs, nil
}
