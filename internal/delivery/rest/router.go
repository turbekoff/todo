package rest

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/turbekoff/todo/internal/delivery/rest/handlers"
	"github.com/turbekoff/todo/internal/delivery/rest/middlewares"
	"github.com/turbekoff/todo/internal/service"
	"golang.org/x/exp/slog"
)

func NewRouter(
	log *slog.Logger,
	userService service.UserService,
	taskService service.TaskService,
	sessionService service.SessionService,
) *chi.Mux {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID)
	router.Use(middlewares.Logger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(sessionService))
		r.Post("/api/v1/logout", handlers.NewLogout(log, sessionService))
		r.Get("/api/v1/profile", handlers.NewProfile(log, userService))
		r.Put("/api/v1/profile", handlers.NewUpdateProfile(log, userService))
		r.Delete("/api/v1/profile", handlers.NewDelete(log, userService))

		r.Post("/api/v1/task", handlers.NewCreateTask(log, taskService))
		r.Get("/api/v1/task", handlers.NewReadTasks(log, taskService))
		r.Get("/api/v1/task/{id}", handlers.NewReadTask(log, taskService))
		r.Put("/api/v1/task/{id}", handlers.NewUpdateTask(log, taskService))
		r.Delete("/api/v1/task/{id}", handlers.NewDeleteTask(log, taskService))
	})

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/signup", handlers.NewSignup(log, userService))
		r.Post("/api/v1/signin", handlers.NewSignin(log, sessionService))
	})

	return router
}
