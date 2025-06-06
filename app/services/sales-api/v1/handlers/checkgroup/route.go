package checkgroup

import (
	"net/http"

	"github.com/1core-dev/go-service/pkg/logger"
	"github.com/1core-dev/go-service/pkg/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	handler := New(cfg.Build, cfg.Log, cfg.DB)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", handler.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", handler.Liveness)
}
