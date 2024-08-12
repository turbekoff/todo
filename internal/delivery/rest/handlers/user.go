package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/turbekoff/todo/internal/delivery/rest/dto"
	"github.com/turbekoff/todo/internal/service"
	"golang.org/x/exp/slog"
)

func NewSignup(log *slog.Logger, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		bind := &dto.UserRequest{}
		if err := render.Bind(r, bind); err != nil {
			log.Error("failed to load request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		err := userService.Create(bind.Name, bind.Password)
		if err != nil {
			log.Error("failed to create user", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		render.Status(r, http.StatusCreated)
	}
}

func NewProfile(log *slog.Logger, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		user, err := userService.Read(fmt.Sprint(r.Context().Value("auth.id")))
		if err != nil {
			log.Error("failed to load request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusInternalServerError, Err: err.Error()})
			return
		}

		render.Render(w, r, &dto.UserResponce{
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}

func NewUpdateProfile(log *slog.Logger, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		bind := &dto.UserRequest{}
		if err := render.Bind(r, bind); err != nil {
			log.Error("failed to load request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		user, err := userService.Update(fmt.Sprint(r.Context().Value("auth.id")), bind.Name, bind.Password)
		if err != nil {
			log.Error("failed to update user", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		render.Render(w, r, &dto.UserResponce{
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}

func NewDelete(log *slog.Logger, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		err := userService.Delete(fmt.Sprint(r.Context().Value("auth.id")))
		if err != nil {
			log.Error("failed to delete user", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		render.Status(r, http.StatusOK)
	}
}
