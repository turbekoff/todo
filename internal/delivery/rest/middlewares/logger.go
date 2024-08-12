package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"golang.org/x/exp/slog"
)

func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "delivery/rest/middleware/logger"),
		)

		handler := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("requestID", middleware.GetReqID(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("device", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
			)

			wrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", wrap.Status()),
					slog.Int("bytes", wrap.BytesWritten()),
					slog.String("duration", time.Since(start).String()),
				)
			}()

			next.ServeHTTP(wrap, r)
		}

		return http.HandlerFunc(handler)
	}
}
