package blog

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/utility"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ListNotes gets all the Notes from the db then encodes them in a response client
func ListNotes(ctx context.Context, db *mongo.Collection) ([]Note, error) {

	list := []Note{}

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting cursor from note collection.")
	}

	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving note list")
	}
	return list, nil
}

// CreateNote takes data from the client to create a note in the db
func CreateNote(ctx context.Context, db *mongo.Collection, user auth.Claims, newNote NewNote, now time.Time) (*Note, error) {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	var strSlice []string
	// check if prop is provided
	if newNote.PropertyAssociation != nil {
		strSlice = append(newNote.PropertyAssociation, newNote.PropertyAssociation...)
		strSlice = utility.RemoveDuplicateStringValues(strSlice)
	}

	note := Note{
		NoteText:            newNote.NoteText,
		PropertyAssociation: strSlice,
		CreatedAt:           now.UTC(),
		UpdatedAt:           now.UTC(),
	}

	nResult, err := db.InsertOne(ctx, note)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting note : %v", note)
	}

	fmt.Println("nResult : ", nResult)

	return &note, nil
}
