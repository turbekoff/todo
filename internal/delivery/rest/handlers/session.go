package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/turbekoff/todo/internal/delivery/rest/dto"
	"github.com/turbekoff/todo/internal/service"
	"golang.org/x/exp/slog"
)

func NewSignin(log *slog.Logger, sessionService service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		bind := &dto.CreateSessionRequest{}
		if err := render.Bind(r, bind); err != nil {
			log.Error("failed to load request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		tokens, err := sessionService.Create(r.RemoteAddr, bind.Name, bind.Password)
		if err != nil {
			log.Error("failed to create session", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    tokens.Refresh,
			Path:     "/",
			HttpOnly: true,
			Expires:  tokens.RefreshExpireAt,
		})

		render.Render(w, r, &dto.SessionResponce{AccessToken: tokens.Access, RefreshToken: tokens.Refresh, ExpireAt: tokens.AccessExpireAt})
	}
}

func NewLogout(log *slog.Logger, sessionService service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now(),
			MaxAge:   -1,
		})
	}
}
