package podcast

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast type is a group of related episodes
type Podcast struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Author      string             `bson:"author,omitempty" json:"author,omitempty"`
	Subscribers int                `bson:"subscribers,omitempty,default:0" json:"subscribers,omitempty,default:0`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// NewPodcast type is what's required from client to create a new Podcast
type NewPodcast struct {
	Title       string   `bson:"_id,omitempty" json:"_id,omitempty"`
	Author      string   `bson:"title,omitempty" json:"title,omitempty"`
	Subscribers int      `bson:"subscribers,omitempty,default:0" json:"subscribers,omitempty,default:0`
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   bool     `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
}

// Episode is the video or audio content
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PodcastID   primitive.ObjectID `bson:"podcastID,omitempty" json:"podcastID,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty" json:"duration,omitempty"`
	Spins       int                `bson:"spins,omitempty,default:0" json:"spins,omitempty,default:0`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// NewEpisode is the video or audio content
type NewEpisode struct {
	PodcastID   primitive.ObjectID `bson:"podcastID,omitempty" json:"podcastID,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty" json:"duration,omitempty"`
	Spins       int                `bson:"spins,omitempty,default:0" json:"spins,omitempty,default:0`
	Published   bool               `bson:"published,omitempty,default:false" json:"published,omitempty,default:false"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}
