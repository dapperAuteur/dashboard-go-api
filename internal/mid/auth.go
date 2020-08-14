package mid

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.opencensus.io/trace"
	// "github.com/pkg/errors"
)

// ErrForbidden is returned when an authenticated user does not have a sufficient role for an action.
var ErrForbidden = web.NewRequestError(errors.New("you are NOT authorized for that action"), http.StatusForbidden,)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(authenticator *auth.Authenticator) web.Middleware {
	
	// This is the actual middleware function to be executed.
	f := func (after web.Handler) web.Handler {
		
		// Wrap this handler around the next one provided.
		h := func (ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.mid.Auth")
			defer span.End()
			
			// Parse the authorization header.
			// Expected header is of the format `Bearer <token>`.
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			_, span = trace.StartSpan(ctx, "internal.ParseClaims")
			claims, err := authenticator.ParseClaims(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			span.End()

			// Add claims to the context so they can be retrieved later.
			ctx = context.WithValue(ctx, auth.Key, claims)

			return after(ctx, w, r)
		}

		return h
	}

	return f
}

// HasRole validates that an authenticated user has at least one role from a specified list.
// This method constructs the actual function that is used.
func HasRole(roles ...string) web.Middleware {
	
	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		
		h := func (ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.mid.HasRole")
			defer span.End()
			
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from context: HasRole called without/before Authenticate")
			}

			if !claims.HasRole(roles...) {
				return ErrForbidden
			}

			return after(ctx, w, r)
		}

		return h
	}

	return f
}