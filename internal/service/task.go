package service

import (
	"errors"
	"strings"
	"time"

	"github.com/turbekoff/todo/internal/domain/entities"
	"github.com/turbekoff/todo/internal/domain/repositories"
)

type taskService struct {
	userRepository repositories.UserRepository
	taskRepository repositories.TaskRepository
}

func NewTaskService(
	userRepository repositories.UserRepository,
	taskRepository repositories.TaskRepository,
) TaskService {
	return &taskService{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (s *taskService) Create(owner, name string, completed bool) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("empty name specified")
	}

	_, err := s.userRepository.Read(owner)
	if err != nil {
		return nil
	}

	now := time.Now()
	task := &entities.Task{
		Owner:     owner,
		Name:      name,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.taskRepository.Create(task)
}

func (s *taskService) Read(id string) (*entities.Task, error) {
	return s.taskRepository.Read(id)
}

func (s *taskService) ReadAllByOwner(owner string) ([]*entities.Task, error) {
	return s.taskRepository.ReadAllByOwner(owner)
}

func (s *taskService) Update(id, name string, completed bool) (*entities.Task, error) {
	task, err := s.taskRepository.Read(id)
	if err != nil {
		return nil, err
	}

	task.Name = name
	task.Completed = completed
	task.UpdatedAt = time.Now()

	if err = s.taskRepository.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) Delete(id string) error {
	return s.taskRepository.Delete(id)
}
