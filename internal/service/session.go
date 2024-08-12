package service

import (
	"errors"
	"time"

	"github.com/turbekoff/todo/internal/config"
	"github.com/turbekoff/todo/internal/domain/entities"
	"github.com/turbekoff/todo/internal/domain/repositories"
	"github.com/turbekoff/todo/pkg/hash"
	"github.com/turbekoff/todo/pkg/jwt"
)

type sessionService struct {
	jwt             *jwt.Manager
	hasher          hash.Hasher
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	userRepository    repositories.UserRepository
	sessionRepository repositories.SessionRepository
}

func NewSessionService(
	hasher hash.Hasher,
	userRepository repositories.UserRepository,
	sessionRepository repositories.SessionRepository,
	JWTConfig *config.JWTConfig,
) SessionService {
	return &sessionService{
		jwt:               jwt.NewManager(JWTConfig.SigningKey),
		hasher:            hasher,
		accessTokenTTL:    JWTConfig.AccessTokenTTL,
		refreshTokenTTL:   JWTConfig.RefreshTokenTTL,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (s *sessionService) create(owner, device string) (*Tokens, error) {
	accessExpireAt := time.Now().Add(s.accessTokenTTL)
	access, err := s.jwt.Generate(owner, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refresh, err := s.jwt.GenerateRefresh()
	if err != nil {
		return nil, err
	}

	refreshExpireAt := time.Now().Add(s.refreshTokenTTL)

	session := &entities.Session{
		Owner:    owner,
		Device:   device,
		Token:    refresh,
		ExpireAt: refreshExpireAt,
	}

	if sessions, err := s.sessionRepository.ReadAllByOwner(owner); !errors.Is(err, repositories.ErrSessionNotFound) {
		if len(sessions) > 10 {
			for _, session := range sessions {
				s.sessionRepository.Delete(session.ID)
			}
		}
	}

	if err := s.sessionRepository.Create(session); err != nil {
		return nil, err
	}

	return &Tokens{Access: access, Refresh: refresh, AccessExpireAt: accessExpireAt, RefreshExpireAt: refreshExpireAt}, nil
}

func (s *sessionService) Create(device, name, password string) (*Tokens, error) {
	user, err := s.userRepository.ReadByName(name)
	if err != nil {
		return nil, err
	}

	if err = s.hasher.Compare(password, user.Password); err != nil {
		return nil, err
	}

	return s.create(user.ID, device)
}

func (s *sessionService) Refresh(device string, token string) (*Tokens, error) {
	session, err := s.sessionRepository.ReadByToken(token)
	if err != nil {
		return nil, err
	}

	s.sessionRepository.Delete(session.ID)

	if time.Now().After(session.ExpireAt) {
		return nil, errors.New("refresh token expired")
	}

	if session.Device != device {
		return nil, errors.New("device doesn't match")
	}

	return s.create(session.Owner, device)
}

func (s *sessionService) VerifyAccess(token string) (string, error) {
	return s.jwt.Parse(token)
}
