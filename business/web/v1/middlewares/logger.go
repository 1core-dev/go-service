package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/1core-dev/go-service/pkg/logger"
	"github.com/1core-dev/go-service/pkg/web"
)

func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(ctx, "request started", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return err
		}

		return h
	}

	return m
}
