package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// App is the entrypoint for all web applications; http package for entire project
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp knows how to construct internal state for an App
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger}
}

// Handle connects a method and URL pattern to a particular application handler
func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHttp(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
