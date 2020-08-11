package user

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser inserts a new user into the database.
func CreateUser(ctx context.Context, db *mongo.Collection, n NewUser, now time.Time) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		CreatedAt:    now.UTC(),
		UpdatedAt:    now.UTC(),
	}

	userResult, err := db.InsertOne(ctx, u)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting User: %v", u)
	}

	returnedUser := db.FindOne(ctx, userResult.InsertedID)

	fmt.Printf("returnedUser  %v: ", returnedUser)

	return &u, nil
}
