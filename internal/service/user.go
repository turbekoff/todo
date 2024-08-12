package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/turbekoff/todo/internal/domain/entities"
	"github.com/turbekoff/todo/internal/domain/repositories"
	"github.com/turbekoff/todo/pkg/hash"
)

var userNameExpression = regexp.MustCompile(`^[0-9A-Za-z]{8,30}$`)

type userService struct {
	hasher            hash.Hasher
	userRepository    repositories.UserRepository
	taskRepository    repositories.TaskRepository
	sessionRepository repositories.SessionRepository
}

func NewUserService(
	hasher hash.Hasher,
	userRepository repositories.UserRepository,
	taskRepository repositories.TaskRepository,
	sessionRepository repositories.SessionRepository,
) UserService {
	return &userService{
		hasher:            hasher,
		userRepository:    userRepository,
		taskRepository:    taskRepository,
		sessionRepository: sessionRepository,
	}
}

func (s *userService) validate(name, password string) error {
	if !userNameExpression.MatchString(name) {
		return errors.New("name must consist of 8-30 latin letters and digits")
	}

	if _, err := s.userRepository.ReadByName(name); !errors.Is(err, repositories.ErrUserNotFound) {
		return errors.New("user with specified name already exists")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, letter := range password {
		switch {
		case 'a' <= letter && letter <= 'z':
			hasLower = true
		case 'A' <= letter && letter <= 'Z':
			hasUpper = true
		case '0' <= letter && letter <= '9':
			hasDigit = true
		case strings.ContainsAny(string(letter), "!@#$&*^()_-+\\"):
			hasSpecial = true
		}
	}

	if !hasLower {
		return errors.New("password must contain latin lowercase letters")
	}

	if !hasUpper {
		return errors.New("password must contain latin uppercase letters")
	}

	if !hasDigit {
		return errors.New("password must contain digits")
	}

	if !hasSpecial {
		return errors.New("password must contain special characters: !@#$&*^()_-+\\")
	}

	return nil
}

func (s *userService) Create(name, password string) error {
	if err := s.validate(name, password); err != nil {
		return err
	}

	password, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &entities.User{
		Name:      name,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.userRepository.Create(user)
}

func (s *userService) Read(id string) (*entities.User, error) {
	return s.userRepository.Read(id)
}

func (s *userService) Update(id, name, password string) (*entities.User, error) {
	if err := s.validate(name, password); err != nil {
		return nil, err
	}

	user, err := s.userRepository.Read(id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Password = password
	user.UpdatedAt = time.Now()

	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id string) error {
	err := s.userRepository.Delete(id)
	if err != nil {
		return err
	}

	tasks, _ := s.taskRepository.ReadAllByOwner(id)
	for _, task := range tasks {
		s.taskRepository.Delete(task.ID)
	}

	sessions, _ := s.sessionRepository.ReadAllByOwner(id)
	for _, session := range sessions {
		s.sessionRepository.Delete(session.ID)
	}

	return nil
}
