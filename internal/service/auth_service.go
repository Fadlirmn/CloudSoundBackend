package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sumbul/music-player-backend/internal/auth"
	"github.com/sumbul/music-player-backend/internal/models"
	"github.com/sumbul/music-player-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, string, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(name, email, password string) (*models.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	user := &models.User{
		ID:           id.String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	// Update last seen
	_ = s.repo.UpdateLastSeen(user.ID)

	return token, user, nil
}
