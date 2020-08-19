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

	sliceStrings = append(sliceStrings, sliceStrings...)

	oIDs, err := SliceStringsToObjectIDs(sliceStrings)
	if err != nil {
		return
	}

	// fmt.Println("oIDs", oIDs)
	lenOIDs := len(oIDs)
	fmt.Println("lenOIDs", lenOIDs)

	lenSliceStrings := len(sliceStrings)
	fmt.Println("lenSliceStrings", lenSliceStrings)

	// test removeDuplicateValues()

	// Assigning values to the slice
	intSliceValues := []int{1, 2, 3, 4, 5, 2, 3, 5, 7, 9, 6, 7}

	// Printing original value of slice
	fmt.Println("intSliceValues", intSliceValues)

	// ==
	// Calling function where we are removing the duplicate ints
	removeDuplicateIntValues := RemoveDuplicateIntValues(intSliceValues)

	// Printing the filtered slice without duplicates
	fmt.Println(removeDuplicateIntValues)

	// ==
	// Comparing strings
	removeDuplicateStringValues := RemoveDuplicateStringValues(sliceStrings)

	fmt.Println("removeDuplicateStringValues", removeDuplicateStringValues)
	fmt.Println("len(removeDuplicateStringValues)", len(removeDuplicateStringValues))

	// ==
	// Comparing ObjectIDs
	removeDuplicateObjectIDValues := RemoveDuplicateObjectIDValues(oIDs)

	fmt.Println("removeDuplicateObjectIDValues", removeDuplicateObjectIDValues)
	fmt.Println("len(removeDuplicateObjectIDValues)", len(removeDuplicateObjectIDValues))

}

// SliceStringsToObjectIDs takes a slice of strings and converts them to a slice of ObjectIDs
func SliceStringsToObjectIDs(sliceStrings []string) ([]primitive.ObjectID, error) {

	var sliceObjectIDs []primitive.ObjectID

	if sliceStrings != nil {

		for _, stringID := range sliceStrings {

			// fmt.Println("stringID : \n", stringID)
			objectID, err := primitive.ObjectIDFromHex(stringID)
			if err != nil {
				return nil, err
			}
			sliceObjectIDs = append(sliceObjectIDs, objectID)
		}
	}
	// fmt.Println("sliceObjectIDs", sliceObjectIDs)

	return sliceObjectIDs, nil
}

// RemoveDuplicateIntValues will remove duplicate items from a slice.
// Then return the slice with all unique values.
func RemoveDuplicateIntValues(intSlice []int) []int {

	// Go type map looks like this:
	// 			 map[KeyType]ValueType
	keys := make(map[int]bool)
	list := []int{}

	// If the key(values of the slice) is not equal
	// to the already present value in the new slice (list)
	// then we append it. Else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDuplicateStringValues will remove duplicate items from a slice.
// Then return the slice with all unique values.
func RemoveDuplicateStringValues(stringSlice []string) []string {

	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDuplicateObjectIDValues will remove duplicate items from a slice.
// Then return the slice with all unique values.
func RemoveDuplicateObjectIDValues(objIDSlice []primitive.ObjectID) []primitive.ObjectID {

	keys := make(map[primitive.ObjectID]bool)
	list := []primitive.ObjectID{}

	for _, entry := range objIDSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
