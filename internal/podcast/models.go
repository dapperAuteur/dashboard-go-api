package podcast

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast type is a group of related episodes
type Podcast struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Author    string             `bson:"author,omitempty" json:"author,omitempty"`
	Tags      []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
