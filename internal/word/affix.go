package word

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
