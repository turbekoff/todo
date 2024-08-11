package repositories

import "github.com/turbekoff/todo/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	Read(id string) (*entities.User, error)
	ReadByName(name string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id string) error
}
