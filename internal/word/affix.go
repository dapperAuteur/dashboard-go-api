package word

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

// AffixList gets all the Affixes from the db then encodes them in a response client.
func AffixList(ctx context.Context, db *mongo.Collection) ([]Affix, error) {

	affixList := []Affix{}

	affixCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting affixCursor retrieving affix list")
	}

	if err = affixCursor.All(ctx, &affixList); err != nil {
		return nil, errors.Wrapf(err, "retrieving affix list")
	}

	return affixList, nil
}

// RetrieveAffixByID gets the first Affix in the db with the provided _id
func RetrieveAffixByID(ctx context.Context, db *mongo.Collection, _id string) (*Affix, error) {

	var affix Affix

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&affix); err != nil {
		return nil, apierror.ErrNotFound
	}

	fmt.Println("affix found : ", affix)

	return &affix, nil
}

// CreateAffix adds a Affix to the database.
// It returns the created Affix with the fields populated, NOT the ID field tho'.
func CreateAffix(ctx context.Context, db *mongo.Collection, user auth.Claims, newAffix NewAffix, now time.Time) (*Affix, error) {

	isAdmin := user.HasRole(auth.RoleAdmin)
	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	affix := Affix{
		AffixType: newAffix.AffixType,
		Example:   newAffix.Example,
		Meaning:   newAffix.Meaning,
		Media:     newAffix.Media,
		Morpheme:  newAffix.Morpheme,
		Note:      newAffix.Note,
		Tongue:    newAffix.Tongue,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
	}

	affixResult, err := db.InsertOne(ctx, affix)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Affix: %v", newAffix)
	}

	fmt.Println("affixResult : ", affixResult)

	return &affix, nil
}

// UpdateOneAffix modifies data about one Affix.
// It will ERROR if the specified affixID is invalid or does NOT reference an existing Affix.
func UpdateOneAffix(ctx context.Context, db *mongo.Collection, user auth.Claims, affixID string, updateAffix UpdateAffix, now time.Time) error {

	affixObjectID, err := primitive.ObjectIDFromHex(affixID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundAffix, err := RetrieveAffixByID(ctx, db, affixID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("affix to update found %+v : \n", foundAffix)

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	affix := Affix{}

	if updateAffix.AffixType != nil {
		objectIDs := append(*updateAffix.AffixType, foundAffix.AffixType...)
		uniqueAffixTypeIds := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}
		affix.AffixType = uniqueAffixTypeIds
	}

	if updateAffix.Meaning != nil {
		objectIDs := append(*updateAffix.Meaning, foundAffix.Meaning...)
		uniqueMeaning := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}
		affix.Meaning = uniqueMeaning
	}

	if updateAffix.Example != nil {
		objectIDs := append(*updateAffix.Example, foundAffix.Example...)
		uniqueExample := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}
		affix.Example = uniqueExample
	}

	if updateAffix.Media != nil {
		objectIDs := append(*updateAffix.Media, foundAffix.Media...)
		uniqueMediaIds := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}
		affix.Media = uniqueMediaIds
	}

	if updateAffix.Note != nil {
		objectIDs := append(*updateAffix.Note, foundAffix.Note...)
		uniqueNoteIds := utility.RemoveDuplicateStringValues(objectIDs)
		if err != nil {
			return err
		}
		affix.Note = uniqueNoteIds
	}

	if updateAffix.Morpheme != nil {
		affix.Morpheme = *updateAffix.Morpheme
	}

	if updateAffix.Tongue != nil {
		affix.Tongue = *updateAffix.Tongue
	}

	affix.ID = affixObjectID

	affix.UpdatedAt = now

	updateA := bson.M{
		"$set": affix,
	}

	fmt.Printf("affix changes set %v : \n", updateA)

	affixResult, err := db.UpdateOne(ctx, bson.M{"_id": affixObjectID}, updateA)
	if err != nil {
		return errors.Wrap(err, "updating affix")
	}

	fmt.Printf("affixResult updated %v : \n", affixResult)

	return nil
}

// DeleteAffixByID removes the Affix identified by a given ID.
func DeleteAffixByID(ctx context.Context, db *mongo.Collection, user auth.Claims, affixID string) error {

	affixObjectID, err := primitive.ObjectIDFromHex(affixID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	_, err = RetrieveAffixByID(ctx, db, affixID)
	if err != nil {
		return apierror.ErrNotFound
	}

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	result, err := db.DeleteOne(ctx, bson.M{"_id": affixObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting affix %s", affixID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
