package hackgroup

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// Routes add specific routes for this gropu.
func Routes(mux *httptreemux.ContextMux) {
	mux.Handle(http.MethodGet, "/hack", Hack)
}
