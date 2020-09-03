package word

import (
	"context"
	"fmt"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
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
