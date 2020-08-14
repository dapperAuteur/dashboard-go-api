package handlers

import (
	"context"
	"net/http"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"github.com/dapperAuteur/dashboard-go-api/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"

	"github.com/pkg/errors"
)

// Users holds handlers for dealing with user.
type Users struct {
	DB            *mongo.Collection
	authenticator *auth.Authenticator
}

// Token generates an authentication token for a user.
// The client must include an email and password for the request using HTTP Basic Auth.
// The user will be identified by email and authenticated by their password.
func (u *Users) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "Handlers.Users.Token")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in BasicAuth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := user.Authenticate(ctx, u.DB, v.Start, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			// return errors.Wrap(err, "authenticating")
			return err
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = u.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	// fmt.Printf("%+v\n", claims)

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
