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

// NewPodcast type is what's required from client to create a new Podcast
type NewPodcast struct {
	Title  string   `json:"title,omitempty"`
	Author string   `json:"author,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

// Episode is the video or audio content
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty" json:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty" json:"duration,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// NewEpisode is the video or audio content
type NewEpisode struct {
	Podcast     primitive.ObjectID `bson:"podcast,omitempty" json:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty" json:"duration,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
}
