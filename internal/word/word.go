package word

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WordList gets all the Words form the db then encodes them in a response client.
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

// RetrieveWord gets the first Word in the db with the provided _id

// CreateWord adds a Word to the database.
// It returns the created Word with fields populated, NOT the ID field tho'. FIX LATER.
func CreateWord(ctx context.Context, db *mongo.Collection, user auth.Claims, newWord NewWord, now time.Time) (*Word, error) {

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
