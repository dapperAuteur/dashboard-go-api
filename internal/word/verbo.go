package word

import (
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// RetrieveVerboByID gets the first Verbo in the db with the provided ID.
func RetrieveVerboByID(ctx context.Context, db *mongo.Collection, _id string) (*Verbo, error) {

	var verbo Verbo

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&verbo); err != nil {
		return nil, apierror.ErrNotFound
	}

	fmt.Println("verbo found : ", verbo)

	return &verbo, nil
}

// CreateVerbo adds a Verbo to the database.
// It returns the created Verbo with the fields populated, NOT the ID field tho'.
func CreateVerbo(ctx context.Context, db *mongo.Collection, user auth.Claims, newVerbo NewVerbo, now time.Time) (*Verbo, error) {

	isAdmin := user.HasRole(auth.RoleAdmin)
	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	verbo := Verbo{
		English:              newVerbo.English,
		Reflexive:            newVerbo.Reflexive,
		Irregular:            newVerbo.Irregular,
		CategoriaDeIrregular: newVerbo.CategoriaDeIrregular,
		CambiarDeIrregular:   newVerbo.CambiarDeIrregular,
		Terminacion:          newVerbo.Terminacion,
		Grupo:                newVerbo.Grupo,
		Spanish:              newVerbo.Spanish,
		CreatedAt:            now.UTC(),
		UpdatedAt:            now.UTC(),
	}

	verboResult, err := db.InsertOne(ctx, verbo)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Verbo: %v", newVerbo)
	}

	fmt.Println("verboResult : ", verboResult)

	return &verbo, nil
}
