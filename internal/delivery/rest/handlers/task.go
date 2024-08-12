package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/turbekoff/todo/internal/delivery/rest/dto"
	"github.com/turbekoff/todo/internal/domain/repositories"
	"github.com/turbekoff/todo/internal/service"
	"golang.org/x/exp/slog"
)

var ErrTaskAuthorization = errors.New("you don't have authorization to view this task")

func NewCreateTask(log *slog.Logger, taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		bind := &dto.TaskRequest{}
		if err := render.Bind(r, bind); err != nil {
			log.Error("failed to load request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		err := taskService.Create(fmt.Sprint(r.Context().Value("auth.id")), bind.Name, bind.Completed)
		if err != nil {
			log.Error("failed to create task", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		render.Status(r, http.StatusCreated)
	}
}

func NewReadTask(log *slog.Logger, taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		log.Debug(chi.URLParam(r, "id"))

		task, err := taskService.Read(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to read task", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusInternalServerError, Err: err.Error()})
			return
		}

		if task.Owner != fmt.Sprint(r.Context().Value("auth.id")) {
			log.Error("failed to read task", slog.Attr{Key: "error", Value: slog.StringValue(ErrTaskAuthorization.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusForbidden, Err: ErrTaskAuthorization.Error()})
			return
		}

		render.Render(w, r, &dto.TaskResponce{
			ID:        task.ID,
			Name:      task.Name,
			Completed: task.Completed,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		})
	}
}

func NewReadTasks(log *slog.Logger, taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		tasks, err := taskService.ReadAllByOwner(fmt.Sprint(r.Context().Value("auth.id")))
		if err != nil && !errors.Is(err, repositories.ErrSessionNotFound) {
			log.Error("failed to read tasks", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		result := &dto.TaskListResponce{}

		for _, task := range tasks {
			*result = append(*result, dto.TaskResponce{
				ID:        task.ID,
				Name:      task.Name,
				Completed: task.Completed,
				CreatedAt: task.CreatedAt,
				UpdatedAt: task.UpdatedAt,
			})
		}

		render.Render(w, r, result)
	}
}

func NewUpdateTask(log *slog.Logger, taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		task, err := taskService.Read(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to read task", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusInternalServerError, Err: err.Error()})
			return
		}

		if task.Owner != fmt.Sprint(r.Context().Value("auth.id")) {
			log.Error("failed to read task", slog.Attr{Key: "error", Value: slog.StringValue(ErrTaskAuthorization.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusForbidden, Err: ErrTaskAuthorization.Error()})
			return
		}

		bind := &dto.TaskRequest{}
		if err := render.Bind(r, bind); err != nil {
			log.Error("failed to load request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		task, err = taskService.Update(chi.URLParam(r, "id"), bind.Name, bind.Completed)
		if err != nil {
			log.Error("failed to update task", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		render.Render(w, r, &dto.TaskResponce{
			ID:        task.ID,
			Name:      task.Name,
			Completed: task.Completed,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		})
	}
}

func NewDeleteTask(log *slog.Logger, taskService service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "signup"),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		task, err := taskService.Read(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to read task", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusInternalServerError, Err: err.Error()})
			return
		}

		if task.Owner != fmt.Sprint(r.Context().Value("auth.id")) {
			log.Error("failed to read task", slog.Attr{Key: "error", Value: slog.StringValue(ErrTaskAuthorization.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusForbidden, Err: ErrTaskAuthorization.Error()})
			return
		}

		err = taskService.Delete(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to delete task", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Render(w, r, &dto.ErrResponce{Code: http.StatusBadRequest, Err: err.Error()})
			return
		}

		render.Status(r, http.StatusOK)
	}
}
