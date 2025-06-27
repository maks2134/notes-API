package service

import (
	"errors"
	"notes-api/internal/model"
	"notes-api/internal/repository"
	"notes-api/internal/util"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(user *model.User) error {
	// Check if user exists
	if _, err := s.userRepo.GetByUsername(user.Username); err == nil {
		return ErrUserExists
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *model.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return "", ErrInvalidCredentials
	}

	return util.GenerateJWT(user.ID)
}
