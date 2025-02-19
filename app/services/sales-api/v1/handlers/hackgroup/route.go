package hackgroup

import (
	"net/http"

	"github.com/1core-dev/go-service/business/web/v1/auth"
	"github.com/1core-dev/go-service/business/web/v1/middlewares"
	"github.com/1core-dev/go-service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth *auth.Auth
}

// Routes add specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := middlewares.Authenticate(cfg.Auth)
	ruleAdmin := middlewares.Authorize(cfg.Auth, auth.RuleAdminOnly)

	app.Handle(http.MethodGet, version, "/hack", Hack)
	app.Handle(http.MethodGet, version, "/hackauth", Hack, authen, ruleAdmin)
}
