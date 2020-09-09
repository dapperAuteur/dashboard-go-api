package user

import (
	"context"
	"fmt"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/apierror"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrAuthenticationFailure occurs when a user attempts to authenticate but anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication Failed")
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

// Authenticate finds a user by their email and verifies their password.
// On success it returns a Claims value representing this user.
// The claims can be used to generate a token for future authentication.
func Authenticate(ctx context.Context, db *mongo.Collection, now time.Time, email, password string) (auth.Claims, error) {

	var u User
	if err := db.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {

		// Normally we would return ErrNotFound in this scenario
		// but we do NOT want to an unauthenticated user which emails are in the system.
		return auth.Claims{}, ErrAuthenticationFailure
	}
	fmt.Println("found user", u)

	// Compare the provided password with the saved hash.
	// Use the bcrypt comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid.
	// Create some claims for the user and generate their token.
	claims := auth.NewClaims(u.ID, u.Roles, now, time.Hour)
	return claims, nil
}

// ListUsers gets all the Users from the database then encodes them in a response client.
func ListUsers(ctx context.Context, db *mongo.Collection, user auth.Claims) ([]User, error) {

	isUser := user.HasRole(auth.RoleUser)
	if !isUser {
		return nil, apierror.ErrForbidden
	}

	list := []User{}

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrapf(err, "getting cursor from User collection when retrieving user list")
	}

	if err = cursor.All(ctx, &list); err != nil {
		return nil, errors.Wrapf(err, "retrieving user list")
	}

	return list, nil
}
