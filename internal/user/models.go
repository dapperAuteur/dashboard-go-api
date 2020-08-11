package user

import (
	"time"
)

// User represents someone with access to our system.
type User struct {
	ID           string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string    `bson:"name,omitempty" json:"name,omitempty"`
	Email        string    `bson:"email,omitempty" json:"email,omitempty"`
	Roles        []string  `bson:"roles,omitempty" json:"roles,omitempty"`
	PasswordHash []byte    `bson:"password_hash" json:"-"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}
