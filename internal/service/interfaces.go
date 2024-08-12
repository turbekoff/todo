package service

import (
	"time"

	"github.com/turbekoff/todo/internal/domain/entities"
)

type Tokens struct {
	Access          string
	Refresh         string
	AccessExpireAt  time.Time
	RefreshExpireAt time.Time
}

type UserService interface {
	Create(name, password string) error
	Read(id string) (*entities.User, error)
	Update(id, name, password string) (*entities.User, error)
	Delete(id string) error
}

type SessionService interface {
	Create(device, name, password string) (*Tokens, error)
	Refresh(device string, token string) (*Tokens, error)
	VerifyAccess(token string) (string, error)
}

type TaskService interface {
	Create(owner, name string, completed bool) error
	Read(id string) (*entities.Task, error)
	ReadAllByOwner(owner string) ([]*entities.Task, error)
	Update(id, name string, completed bool) (*entities.Task, error)
	Delete(id string) error
}
