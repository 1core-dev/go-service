package checkgroup

import (
	"net/http"

	"github.com/1core-dev/go-service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	handler := New(cfg.Build)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", handler.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", handler.Liveness)
}
