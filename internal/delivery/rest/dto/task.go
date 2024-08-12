package dto

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type TaskRequest struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func (task *TaskRequest) Bind(r *http.Request) error {
	return nil
}

type TaskResponce struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (task *TaskResponce) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

type TaskListResponce []TaskResponce

func (tasks *TaskListResponce) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
