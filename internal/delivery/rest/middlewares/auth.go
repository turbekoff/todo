package middlewares

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/turbekoff/todo/internal/delivery/rest/dto"
	"github.com/turbekoff/todo/internal/service"
	"golang.org/x/net/context"
)

func Auth(sessionService service.SessionService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			if len(auth) < 7 {
				render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: "invalid bearer authorization header"})
				return
			}

			authType := strings.ToLower(auth[:6])
			if authType != "bearer" {
				render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: "invalid bearer authorization header"})
				return
			}

			id, err := sessionService.VerifyAccess(strings.TrimSpace(auth[7:]))
			if err != nil {
				render.Render(w, r, &dto.ErrResponce{Code: http.StatusUnauthorized, Err: err.Error()})
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "auth.id", id)
			ctx = context.WithValue(ctx, "auth.device", r.RemoteAddr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
