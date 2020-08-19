package blog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note type is a prop added to other models to add extra info about the element
type Note struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	PropertyAssociation []string           `bson:"prop_assoc,omitempty" json:"prop_assoc,omitempty"`
	NoteText            string             `bson:"note_text,omitempty" json:"note_text,omitempty"`
	CreatedAt           time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt           time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// NewNote type is what's required from the client to create a new Note
type NewNote struct {
	PropertyAssociation []string `bson:"prop_assoc,omitempty" json:"prop_assoc,omitempty"`
	NoteText            string   `bson:"note_text,omitempty" json:"note_text,omitempty"`
}

// UpdateNote defines what information may be provided to modify an existing Note.
// All fields are optional so clients can send just the fields they want changed.
// It uses pointer fields so we can differentiate between a field that was not provided and a field that was provided as explicitly blank.
// Normally we do not want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateNote struct {
	ID                  *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	PropertyAssociation *[]string           `bson:"prop_assoc,omitempty" json:"prop_assoc,omitempty"`
	NoteText            *string             `bson:"note_text,omitempty" json:"note_text,omitempty"`
}
