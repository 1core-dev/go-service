package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
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

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// HandleNoMiddleware sets a handler function for a given HTTP method and path pair
// to the application server mux. Does not include the application middleware.
func (a *App) HandleNoMiddleware(method string, group string, path string, handler Handler) {
	a.handle(method, group, path, handler)
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	a.handle(method, group, path, handler)
}

// handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) handle(method string, group string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := SetValues(r.Context(), &v)
		if err := handler(ctx, w, r); err != nil {
			if validateShutdown(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	a.ContextMux.Handle(method, finalPath, h)
}

// validateShutdown validates the error for special condition that don't
// warrant an actual shutdown by the system
func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection
		// after the RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream
		// where the other end (the peer) has closed the underlying connection.
		// The first write to the closed connection cause the peer to reply with
		// an RST packet indicating that the connection should be terminated immediately.
		// The second write to the socket that has already received the RST
		// causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from
		// the connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the
		// other end (peer) has unexpectedly closed the connection.
		// It happens when you send a packet from your end, but the other end
		// crashes and forcibly closes the connection with the RST packet
		// instead of the TCP FIN, which is used to close a connection under
		// normal circumstances.
		return false

	}

	return true
}
