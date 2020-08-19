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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// RetrieveNote finds a single note by _id
func RetrieveNote(ctx context.Context, db *mongo.Collection, _id string) (*Note, error) {

	var note Note

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&note); err != nil {
		return nil, apierror.ErrNotFound
	}

	return &note, nil
}

// UpdateOneNote modifies data about a note.
// It will error if the specified _id is invalid or does NOT reference an existing note.
func UpdateOneNote(ctx context.Context, db *mongo.Collection, user auth.Claims, nID string, updateNote UpdateNote, now time.Time) error {

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	foundNote, err := RetrieveNote(ctx, db, nID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("note to update found %+v : \n", foundNote)

	nObjectID, err := primitive.ObjectIDFromHex(nID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	note := Note{}

	if updateNote.PropertyAssociation != nil {
		// take *updateNote.PropertyAssociation
		// confirm unique values
		strSlice := utility.RemoveDuplicateStringValues(*updateNote.PropertyAssociation)
		note.PropertyAssociation = strSlice
	}

	if updateNote.NoteText != nil {
		note.NoteText = *updateNote.NoteText
	}

	note.ID = nObjectID

	note.UpdatedAt = now

	updateN := bson.M{
		"$set": note,
	}

	nResult, err := db.UpdateOne(ctx, bson.M{"_id": nObjectID}, updateN)
	if err != nil {
		return errors.Wrap(err, "updating note")
	}

	fmt.Printf("nResult updated %v : \n", nResult)

	return nil
}

// DeleteNote removes the note identified by a given _id
func DeleteNote(ctx context.Context, db *mongo.Collection, user auth.Claims, noteID string) error {

	var isAdmin = user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	nObjectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundNote, err := RetrieveNote(ctx, db, noteID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("note to delelete found %+v : \n", foundNote)

	result, err := db.DeleteOne(ctx, bson.M{"_id": nObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting note %s", noteID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
