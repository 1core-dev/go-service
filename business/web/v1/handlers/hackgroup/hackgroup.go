package hackgroup

import (
	"context"
	"net/http"

	"github.com/1core-dev/go-service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, r, status, http.StatusOK)
}
