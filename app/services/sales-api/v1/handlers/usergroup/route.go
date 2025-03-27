package usergroup

import (
	"net/http"

	"github.com/1core-dev/go-service/business/core/user"
	"github.com/1core-dev/go-service/business/core/user/stores/userdb"
	"github.com/1core-dev/go-service/business/web/v1/auth"
	"github.com/1core-dev/go-service/business/web/v1/middlewares"
	"github.com/1core-dev/go-service/foundation/logger"
	"github.com/1core-dev/go-service/foundation/web"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authentication := middlewares.Authenticate(cfg.Auth)
	ruleAdmin := middlewares.Authorize(cfg.Auth, auth.RuleAdminOnly)

	usrCore := user.NewCore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))

	handler := New(usrCore, cfg.Auth)
	app.Handle(http.MethodPost, version, "/users", handler.Create)
	app.Handle(http.MethodPost, version, "/usersauth", handler.Create, authentication, ruleAdmin)
}
