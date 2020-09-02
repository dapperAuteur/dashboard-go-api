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

// WordList gets all the Words from the db then encodes them in a response client.
func WordList(ctx context.Context, db *mongo.Collection) ([]Word, error) {

	wordList := []Word{}

	wordCursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting wordCursor retrieving word list")
	}

	if err = wordCursor.All(ctx, &wordList); err != nil {
		return nil, errors.Wrapf(err, "retrieving word list")
	}

	return wordList, nil
}

// RetrieveWordByID gets the first Word in the db with the provided _id
func RetrieveWordByID(ctx context.Context, db *mongo.Collection, _id string) (*Word, error) {

	var word Word

	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	if err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&word); err != nil {
		return nil, apierror.ErrNotFound
	}

	fmt.Println("word found : ", word)

	return &word, nil
}

// RetrieveWord gets the first Word in the db with the provided wd string
func RetrieveWord(ctx context.Context, db *mongo.Collection, wd string) (*Word, error) {

	var word Word

	// Create filter to find word by word parameter
	filter := Word{
		Word: wd,
	}

	if err := db.FindOne(ctx, filter).Decode(&word); err != nil {
		return nil, apierror.ErrNotFound
	}

	fmt.Println("word found : ", word)

	return &word, nil
}

// CreateWord adds a Word to the database.
// It returns the created Word with fields populated, NOT the ID field tho'. FIX LATER.
func CreateWord(ctx context.Context, db *mongo.Collection, user auth.Claims, newWord NewWord, now time.Time) (*Word, error) {

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return nil, apierror.ErrForbidden
	}

	word := Word{
		Meaning:          newWord.Meaning,
		Tongue:           newWord.Tongue,
		InGame:           newWord.InGame,
		IsFourLetterWord: newWord.IsFourLetterWord,
		Word:             newWord.Word,
		FPoints:          newWord.FPoints,
		SPoints:          newWord.SPoints,
		Tier:             newWord.Tier,
		CreatedAt:        now.UTC(),
		UpdatedAt:        now.UTC(),
	}

	wordResult, err := db.InsertOne(ctx, word)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting Word: %v", newWord)
	}
	fmt.Println("wordResult : ", wordResult)

	return &word, nil
}

// UpdateOneWord modifies data about one Word.
// It will ERROR if the specified wordID is invalid or does NOT reference an existing Word.
func UpdateOneWord(ctx context.Context, db *mongo.Collection, user auth.Claims, wordID string, updateWord UpdateWord, now time.Time) error {

	wordObjectID, err := primitive.ObjectIDFromHex(wordID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	foundWord, err := RetrieveWordByID(ctx, db, wordID)
	if err != nil {
		return apierror.ErrNotFound
	}

	fmt.Printf("word to update found %+v : \n", foundWord)

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	word := Word{}

	if updateWord.Meaning != nil {
		word.Meaning = *updateWord.Meaning
	}

	if updateWord.Tongue != nil {
		word.Tongue = *updateWord.Tongue
	}

	if updateWord.Tier != nil {
		word.Tier = *updateWord.Tier
	}

	if updateWord.Word != nil {
		word.Word = *updateWord.Word
	}

	if updateWord.InGame != nil {
		word.InGame = *updateWord.InGame
	}

	if updateWord.SPoints != nil {
		word.SPoints = *updateWord.SPoints
	}

	if updateWord.FPoints != nil {
		word.FPoints = *updateWord.FPoints
	}

	if updateWord.IsFourLetterWord != nil {
		word.IsFourLetterWord = *updateWord.IsFourLetterWord
	}

	word.ID = wordObjectID

	word.UpdatedAt = now

	updateW := bson.M{
		"$set": word,
	}

	fmt.Printf("word changes set %v : \n", updateW)

	wordResult, err := db.UpdateOne(ctx, bson.M{"_id": wordObjectID}, updateW)
	if err != nil {
		return errors.Wrap(err, "updating word")
	}

	fmt.Printf("wordResult updated %v : \n", wordResult)

	return nil
}

// DeleteWord removes the Word identified by a given ID
func DeleteWord(ctx context.Context, db *mongo.Collection, user auth.Claims, wordID string) error {

	fmt.Printf("************** looking for wordID %s : ", wordID)

	wordObjectID, err := primitive.ObjectIDFromHex(wordID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	fmt.Printf("************** looking for wordID %v : ", wordObjectID)

	_, err = RetrieveWordByID(ctx, db, wordID)
	if err != nil {
		return apierror.ErrNotFound
	}

	isAdmin := user.HasRole(auth.RoleAdmin)

	if !isAdmin {
		return apierror.ErrForbidden
	}

	result, err := db.DeleteOne(ctx, bson.M{"_id": wordObjectID})
	if err != nil {
		return errors.Wrapf(err, "deleting word %s", wordID)
	}

	fmt.Print("result of deleting : ", result)

	return nil
}
