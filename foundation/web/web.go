package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// Handler is a type that handles a http request within own "mini framework".
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// NewApp created an App value that handle as et og routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := SetValues(r.Context(), &v)
		if err := handler(ctx, w, r); err != nil {
			return
		}

	}
	a.ContextMux.Handle(method, path, h)
}
