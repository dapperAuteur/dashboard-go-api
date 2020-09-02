package word

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Word type is a group of English words
type Word struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	Meaning          []string           `bson:"meaning,omitempty" json:"meaning,omitempty"`
	Tongue           string             `bson:"tongue,omitempty" json:"tongue,omitempty"`
	Tier             int                `bson:"tier,omitempty" json:"tier,omitempty"`
	Word             string             `bson:"word,omitempty" json:"word,omitempty"`
	InGame           bool               `bson:"in_game,omitempty" json:"in_game,omitempty" default:"false"`
	SPoints          int                `bson:"s_points,omitempty" json:"s_points,omitempty"`
	FPoints          int                `bson:"f_points,omitempty" json:"f_points,omitempty"`
	IsFourLetterWord bool               `bson:"is_four_letter_word,omitempty" json:"is_four_letter_word,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"datetime"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"datetime"`
}

// NewWord type is what's required from the client to create a new Word
type NewWord struct {
	Meaning          []string `bson:"meaning,omitempty" json:"meaning,omitempty"`
	Tongue           string   `bson:"tongue,omitempty" json:"tongue,omitempty"`
	Tier             int      `bson:"tier,omitempty" json:"tier,omitempty"`
	Word             string   `bson:"word,omitempty" json:"word,omitempty"`
	InGame           bool     `bson:"in_game,omitempty" json:"in_game,omitempty" default:"false"`
	SPoints          int      `bson:"s_points,omitempty" json:"s_points,omitempty"`
	FPoints          int      `bson:"f_points,omitempty" json:"f_points,omitempty"`
	IsFourLetterWord bool     `bson:"is_four_letter_word,omitempty" json:"is_four_letter_word,omitempty"`
}

// UpdateWord defines what information may be provided to modify an existing Word.
// All fields are optional so the clients can send just the fields they want changed.
// It uses pointer fields so we can differentiate between a field that was NOT provided and a field that was provided as explicitly blank.
// Normally we do NOT want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateWord struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Meaning          *[]string          `bson:"meaning,omitempty" json:"meaning,omitempty"`
	Tongue           *string            `bson:"tongue,omitempty" json:"tongue,omitempty"`
	Tier             *int               `bson:"tier,omitempty" json:"tier,omitempty"`
	Word             *string            `bson:"word,omitempty" json:"word,omitempty"`
	InGame           *bool              `bson:"in_game,omitempty" json:"in_game,omitempty" default:"false"`
	SPoints          *int               `bson:"s_points,omitempty" json:"s_points,omitempty"`
	FPoints          *int               `bson:"f_points,omitempty" json:"f_points,omitempty"`
	IsFourLetterWord *bool              `bson:"is_four_letter_word,omitempty" json:"is_four_letter_word,omitempty"`
}

// Affix type is a group of related Affixes
type Affix struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	Morpheme  string             `bson:"morpheme,omitempty" json:"morpheme,omitempty"`
	Meaning   []string           `bson:"meaning,omitempty" json:"meaning,omitempty"`
	Tongue    string             `bson:"tongue,omitempty" json:"tongue,omitempty"`
	Example   []string           `bson:"example,omitempty" json:"example,omitempty"`
	AffixType []string           `bson:"affix_type,omitempty" json:"affix_type,omitempty"`
	Media     []string           `bson:"media,omitempty" json:"media,omitempty"`
	Note      []string           `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"datetime"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"datetime"`
}

// NewAffix type is what's required from the client to create a new Affix.
type NewAffix struct {
	Morpheme  string   `bson:"morpheme,omitempty" json:"morpheme,omitempty"`
	Meaning   []string `bson:"meaning,omitempty" json:"meaning,omitempty"`
	Tongue    string   `bson:"tongue,omitempty" json:"tongue,omitempty"`
	Example   []string `bson:"example,omitempty" json:"example,omitempty"`
	AffixType []string `bson:"affix_type,omitempty" json:"affix_type,omitempty"`
	Media     []string `bson:"media,omitempty" json:"media,omitempty"`
	Note      []string `bson:"note,omitempty" json:"note,omitempty"`
}

// UpdateAffix defines what information may be provided to modify an existing Word.
// All fields are optional so clients can send just the fields they want changed.
// It uses pointer fields so we can differentiate between a field that was NOT provided and a field that was provided as explicitly blank.
// Normally we do NOT want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateAffix struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	Morpheme  *string            `bson:"morpheme,omitempty" json:"morpheme,omitempty"`
	Meaning   *[]string          `bson:"meaning,omitempty" json:"meaning,omitempty"`
	Tongue    *string            `bson:"tongue,omitempty" json:"tongue,omitempty"`
	Example   *[]string          `bson:"example,omitempty" json:"example,omitempty"`
	AffixType *[]string          `bson:"affix_type,omitempty" json:"affix_type,omitempty"`
	Media     *[]string          `bson:"media,omitempty" json:"media,omitempty"`
	Note      *[]string          `bson:"note,omitempty" json:"note,omitempty"`
}
