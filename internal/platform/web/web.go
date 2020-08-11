package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// Handler is the signature that all application handlers will implement.
type Handler func(http.ResponseWriter, *http.Request) error

// App is the entrypoint for all web applications; http package for entire project
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
}

// NewApp knows how to construct internal state for an App
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
		mw:  mw,
	}
}

// Handle connects a method and URL pattern to a particular application handler
func (a *App) Handle(method, pattern string, h Handler) {

	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			Start: time.Now(),
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyValues, &v)
		r = r.WithContext(ctx)

		if err := h(w, r); err != nil {
			a.log.Printf("ERROR : Unhandled error %v", err)

		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
