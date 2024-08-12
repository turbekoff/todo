package dto

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type CreateSessionRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (session *CreateSessionRequest) Bind(r *http.Request) error {
	return nil
}

type SessionResponce struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpireAt     time.Time `json:"expireAt"`
}

func (session *SessionResponce) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
