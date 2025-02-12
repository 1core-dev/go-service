package middlewares

import (
	"context"
	"net/http"

	"github.com/1core-dev/go-service/business/web/v1/auth"
	"github.com/1core-dev/go-service/business/web/v1/response"
	"github.com/1core-dev/go-service/foundation/logger"
	"github.com/1core-dev/go-service/foundation/web"
)

// Errors handle errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				var e response.ErrorDocument
				var status int

				switch {
				case response.IsError(err):
					reqErr := response.GetError(err)
					e = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					e = response.ErrorDocument{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					e = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, e, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it back
				// to the base handler to shutdown the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}
		return h
	}
	return m
}
