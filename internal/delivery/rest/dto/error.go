package dto

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponce struct {
	Code int    `json:"code"`
	Err  string `json:"message"`
}

func (err *ErrResponce) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, err.Code)
	return nil
}
