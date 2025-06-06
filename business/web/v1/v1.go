package v1

import (
	"os"

	"github.com/1core-dev/go-service/business/web/v1/auth"
	"github.com/1core-dev/go-service/business/web/v1/middlewares"
	"github.com/1core-dev/go-service/pkg/logger"
	"github.com/1core-dev/go-service/pkg/web"
	"github.com/jmoiron/sqlx"
)

// APIMuxConfig contains all the mandatory system required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	Auth     *auth.Auth
	DB       *sqlx.DB
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg APIMuxConfig)
}

// APIMux constructs a http.Handler will all application routes defined.
func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		middlewares.Logger(cfg.Log),
		middlewares.Errors(cfg.Log),
		middlewares.Metrics(),
		middlewares.Panics(),
	)

	routeAdder.Add(app, cfg)

	return app
}
