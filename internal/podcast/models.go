package podcast

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast type is a group of related episodes
type Podcast struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty" validate:"required"`
	Author      string             `bson:"author,omitempty" json:"author,omitempty" validate:"required"`
	Subscribers int                `bson:"subscribers,omitempty,default:0" json:"subscribers,omitempty,default:0" validate:"gte=0"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// NewPodcast type is what's required from client to create a new Podcast
type NewPodcast struct {
	Title       string   `bson:"title,omitempty" json:"title,omitempty" validate:"required"`
	Author      string   `bson:"author,omitempty" json:"author,omitempty" validate:"required"`
	Subscribers int      `bson:"subscribers,omitempty,default:0" json:"subscribers,omitempty,default:0" validate:"gte=0"`
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   bool     `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
}

// UpdatePodcast defines what information may be provided to modify an
// existing Podcast. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that was not
// provided and a field that was provided as explicitly blank. Normally we do not want to
// use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdatePodcast struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       *string            `bson:"title,omitempty" json:"title,omitempty" validate:"required"`
	Author      *string            `bson:"author,omitempty" json:"author,omitempty" validate:"required"`
	Subscribers *int               `bson:"subscribers,omitempty,default:0" json:"subscribers,omitempty,default:0" validate:"gte=0"`
	Tags        *[]string          `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   *bool              `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
}

// Episode is the video or audio content
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PodcastID   primitive.ObjectID `bson:"podcastID,omitempty" json:"podcastID,omitempty" validate:"required"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty" validate:"required"`
	Description string             `bson:"description,omitempty" json:"description,omitempty" validate:"required"`
	Duration    int32              `bson:"duration,omitempty" json:"duration,omitempty" validate:"gte=0"`
	Spins       int                `bson:"spins,omitempty,default:0" json:"spins,omitempty,default:0" validate:"gte=0"`
	Published   bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// NewEpisode is the video or audio content
type NewEpisode struct {
	PodcastID   primitive.ObjectID `bson:"podcastID,omitempty" json:"podcastID,omitempty" validate:"required"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty" validate:"required"`
	Description string             `bson:"description,omitempty" json:"description,omitempty" validate:"required"`
	Duration    int32              `bson:"duration,omitempty" json:"duration,omitempty" validate:"gte=0"`
	Spins       int                `bson:"spins,omitempty,default:0" json:"spins,omitempty,default:0" validate:"gte=0"`
	Published   bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}

// UpdateEpisode defines what information may be provided to modify an
// existing Episode. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that was not
// provided and a field that was provided as explicitly blank. Normally we do not want to
// use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateEpisode struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	PodcastID   *primitive.ObjectID `bson:"podcastID,omitempty" json:"podcastID,omitempty" validate:"required"`
	Title       *string             `bson:"title,omitempty" json:"title,omitempty" validate:"required"`
	Description *string             `bson:"description,omitempty" json:"description,omitempty" validate:"required"`
	Duration    *int32              `bson:"duration,omitempty" json:"duration,omitempty" validate:"gte=0"`
	Spins       *int                `bson:"spins,omitempty,default:0" json:"spins,omitempty,default:0" validate:"gte=0"`
	Published   *bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	Tags        *[]string           `bson:"tags,omitempty" json:"tags,omitempty"`
}
