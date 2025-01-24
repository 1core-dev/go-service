package hackgroup

import (
	"net/http"

	"github.com/1core-dev/go-service/foundation/web"
)

// Routes add specific routes for this gropu.
func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
