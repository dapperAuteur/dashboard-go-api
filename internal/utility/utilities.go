package utility

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
