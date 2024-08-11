package repositories

import (
	"github.com/turbekoff/todo/internal/domain/entities"
)

type SessionRepository interface {
	Create(session *entities.Session) error
	Read(id string) (*entities.Session, error)
	ReadByToken(token string) (*entities.Session, error)
	ReadAllByOwner(owner string) ([]*entities.Session, error)
	Update(session *entities.Session) error
	Delete(id string) error
}
