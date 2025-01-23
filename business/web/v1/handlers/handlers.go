package handlers

import (
	v1 "github.com/1core-dev/go-service/business/web/v1"
	"github.com/1core-dev/go-service/business/web/v1/handlers/hackgroup"
	"github.com/dimfeld/httptreemux/v5"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(mux *httptreemux.ContextMux, cfg v1.APIMuxConfig) {
	hackgroup.Routes(mux)
}
