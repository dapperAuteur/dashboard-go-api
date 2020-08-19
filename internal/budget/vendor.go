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

// ListVendors gets all the Vendors from the db then encodes them in a response client
func ListVendors(ctx context.Context, db *mongo.Collection) ([]Vendor, error) {

	list := []Vendor{}

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting cursor from vendor collection.")
	}

	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving vendor list")
	}

	return list, nil
}

// CreateVendor takes data from the client to create a vendor in the db
func CreateVendor(ctx context.Context, db *mongo.Collection, user auth.Claims, newVendor NewVendor, now time.Time) (*Vendor, error) {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	var tranxObjectIDs []primitive.ObjectID

	// check if prop is provided
	if newVendor.TransactionIDs != nil {
		// convert []newVendor.TransactionIDs (ObjectID) to []string
		objIDs, err := utility.SliceStringsToObjectIDs(*newVendor.TransactionIDs)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, objIDs...)
		tranxObjectIDs = utility.RemoveDuplicateObjectIDValues(objIDs)
	}

	vendor := Vendor{
		VendorName:     newVendor.VendorName,
		TransactionIDs: tranxObjectIDs,
		CreatedAt:      now.UTC(),
		UpdatedAt:      now.UTC(),
	}

	vResult, err := db.InsertOne(ctx, vendor)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting vendor : %v", vendor)
	}

	fmt.Println("vResult : ", vResult)

	return &vendor, nil
}

// RetrieveVendor finds a single vendor by _id
func RetrieveVendor(ctx context.Context, db *mongo.Collection, _id string) (*Vendor, error) {

	var vendor Vendor

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&vendor); err != nil {
		return nil, apierror.ErrNotFound
	}

	return &vendor, nil
}

// UpdateOneVendor modifies data about a vendor.
// It will error if the specified _id is invalid or does NOT reference an existing vendor.
func UpdateOneVendor(ctx context.Context, db *mongo.Collection, user auth.Claims, vID string, updateVendor UpdateVendor, now time.Time) error {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	foundVendor, err := RetrieveVendor(ctx, db, vID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("vendor to update found %+v : \n", foundVendor)

	vObjectID, err := primitive.ObjectIDFromHex(vID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	vendor := Vendor{}

	if updateVendor.VendorName != nil {
		vendor.VendorName = *updateVendor.VendorName
	}

	if updateVendor.TransactionIDs != nil {
		// take *updateVendor.TransactionIDs
		// convert to []primitive.ObjectID and return
		vendorObjIDs, err := utility.SliceStringsToObjectIDs(*updateVendor.TransactionIDs)
		if err != nil {
			return err
		}
		objectIDs := append(vendorObjIDs, foundVendor.TransactionIDs...)
		uniqueVendorTranxObjIDs := utility.RemoveDuplicateObjectIDValues(objectIDs)
		vendor.TransactionIDs = uniqueVendorTranxObjIDs
	}

	vendor.ID = vObjectID

	vendor.UpdatedAt = now

	updateV := bson.M{
		"$set": vendor,
	}

	vResult, err := db.UpdateOne(ctx, bson.M{"_id": vObjectID}, updateV)
	if err != nil {
		return errors.Wrap(err, "updating vendor")
	}

	fmt.Printf("vResult updated %v : \n", vResult)

	return nil
}

// DeleteVendor removes the vendor identified by a given _id
func DeleteVendor(ctx context.Context, db *mongo.Collection, user auth.Claims, vendorID string) error {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	vObjectID, err := primitive.ObjectIDFromHex(vendorID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundVendor, err := RetrieveVendor(ctx, db, vendorID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("vendor to delelete found %+v : \n", foundVendor)

	result, err := db.DeleteOne(ctx, bson.M{"_id": vObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting vendor %s", vendorID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
