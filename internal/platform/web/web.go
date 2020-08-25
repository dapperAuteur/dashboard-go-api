package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	TraceID    string
	StatusCode int
	Start      time.Time
}

// Handler is the signature that all application handlers will implement.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// App is the entrypoint for all web applications; http package for entire project
type App struct {
	mux      *chi.Mux
	log      *log.Logger
	mw       []Middleware
	och      *ochttp.Handler
	shutdown chan os.Signal
}

// NewApp constructs an App to handle a set of routes.
// Any Middleware provided will be ran for every request.
func NewApp(shutdown chan os.Signal, logger *log.Logger, mw ...Middleware) *App {
	app := App{
		mux:      chi.NewRouter(),
		log:      logger,
		mw:       mw,
		shutdown: shutdown,
	}

	app.mux.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	// Create an OpenCensus HTTP Handler which wraps the router.
	// This will start the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote parent if a client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	app.och = &ochttp.Handler{
		Handler:     app.mux,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &app
}

// Handle associates a handler function with an HTTP Method and URL pattern.
//
// It converts our custom handler type to the std lib Handler type. It captures
// errors from the handler and serves them to the client in a uniform way.
func (a *App) Handle(method, pattern string, h Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	h = wrapMiddleware(mw, h)

	// Add the application's general middleware to the handler chain.
	h = wrapMiddleware(a.mw, h)

	// Create a function that conforms to the std lib definition of a handler.
	// This is the first thing that will be executed when this route is called.
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		// Create a Values struct to record state for the request. Store the
		// address in the request's context so it is sent down the call chain.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Start:   time.Now(),
		}

		ctx = context.WithValue(ctx, KeyValues, &v) // commenting out this line will cause an integrity shutdown. used to test self-shutdown works when app integrity is compromised

		// Run the handler chain and catch any propagated error.
		if err := h(ctx, w, r); err != nil {
			a.log.Printf("%s : Unhandled error %+v", v.TraceID, err)
			if IsShutdown(err) {
				a.SignalShutdown()
			}

		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.log.Println("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGSTOP
}
