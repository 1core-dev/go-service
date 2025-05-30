package usergroup

import (
	"net/http"

	"github.com/1core-dev/go-service/business/core/user"
	"github.com/1core-dev/go-service/business/core/user/stores/userdb"
	db "github.com/1core-dev/go-service/business/data/dbsql/pgx"
	"github.com/1core-dev/go-service/business/web/v1/auth"
	"github.com/1core-dev/go-service/business/web/v1/middlewares"
	"github.com/1core-dev/go-service/pkg/logger"
	"github.com/1core-dev/go-service/pkg/web"

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
	ruleAdminOrSubject := middlewares.Authorize(cfg.Auth, auth.RuleAdminOrSubject)
	tx := middlewares.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	usrCore := user.NewCore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))

	handler := New(usrCore, cfg.Auth)
	app.Handle(http.MethodPost, version, "/users", handler.Create)
	app.Handle(http.MethodPost, version, "/userstran", handler.CreateWithTran, authentication, ruleAdmin, tx)
	app.Handle(http.MethodPost, version, "/usersauth", handler.Create, authentication, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users", handler.Query, authentication, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users/:user_id", handler.QueryByID, authentication, ruleAdminOrSubject)
}
