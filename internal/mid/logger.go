package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/platform/web"
	"go.opencensus.io/trace"
)

// Logger will log a line for every request.
// Logger writes some information about the request to the logs in the format:
// (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.mid.Logger")
			defer span.End()

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			err := before(ctx, w, r)

			log.Printf("%s : (%d) : %s %s -> %s (%s)",
				v.TraceID,
				v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}
		return h
	}
	return f
}
