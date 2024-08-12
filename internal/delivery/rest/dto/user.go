package dto

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type UserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (user *UserRequest) Bind(r *http.Request) error {
	return nil
}

type UserResponce struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *UserResponce) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
