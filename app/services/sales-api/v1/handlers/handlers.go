package handlers

import (
	"github.com/1core-dev/go-service/app/services/sales-api/v1/handlers/hackgroup"
	v1 "github.com/1core-dev/go-service/business/web/v1"
	"github.com/1core-dev/go-service/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgroup.Routes(app)
}
