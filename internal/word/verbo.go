package word

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

// VerboList gets all the Verbos from the database then encodes them in a response client.
func VerboList(ctx context.Context, db *mongo.Collection) ([]Verbo, error) {

	verboList := []Verbo{}

	verboCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting verboCursor retrieving verbo list")
	}

	if err = verboCursor.All(ctx, &verboList); err != nil {
		return nil, errors.Wrapf(err, "retrieving verbo list")
	}

	return verboList, nil
}
