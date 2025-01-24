package web

import (
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

// App is the entrypoint into our application ans what configures our context
// object for each of our http handlers.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

// NewApp created an App value that handle as et og routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}
