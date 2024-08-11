package repositories

import "github.com/turbekoff/todo/internal/domain/entities"

type TaskRepository interface {
	Create(task *entities.Task) error
	Read(id string) (*entities.Task, error)
	ReadAllByOwner(owner string) ([]*entities.Task, error)
	Update(task *entities.Task) error
	Delete(id string) error
}
