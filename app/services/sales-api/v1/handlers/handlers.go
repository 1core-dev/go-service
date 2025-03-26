package handlers

import (
	"github.com/1core-dev/go-service/app/services/sales-api/v1/handlers/checkgroup"
	"github.com/1core-dev/go-service/app/services/sales-api/v1/handlers/hackgroup"
	"github.com/1core-dev/go-service/app/services/sales-api/v1/handlers/usergroup"
	v1 "github.com/1core-dev/go-service/business/web/v1"
	"github.com/1core-dev/go-service/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface to add all routes.
func (Routes) Add(app *web.App, apiCfg v1.APIMuxConfig) {
	hackgroup.Routes(app, hackgroup.Config{
		Auth: apiCfg.Auth,
	})

	checkgroup.Routes(app, checkgroup.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
		DB:    apiCfg.DB,
	})

	usergroup.Routes(app, usergroup.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
		DB:    apiCfg.DB,
		Auth:  apiCfg.Auth,
	})
}
