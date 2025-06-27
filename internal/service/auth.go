package service

import (
	"errors"
	"fmt"
	"notes-api/internal/model"
	"notes-api/internal/repository"
	"notes-api/internal/util"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(user *model.User) error {
	// Логируем перед хешированием
	fmt.Printf("Registering user: %s with password: %s\n", user.Username, user.Password)

	hashed, err := util.HashPassword(user.Password)
	if err != nil {
		fmt.Printf("Hashing error: %v\n", err)
		return err
	}

	user.Password = hashed
	fmt.Printf("Storing hash: %s\n", user.Password)

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *model.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("user not found: %v", err) // Добавим логирование
	}

	fmt.Printf("Input password: %s\n", req.Password)
	fmt.Printf("Stored hash: %s\n", user.Password)

	match := util.CheckPasswordHash(req.Password, user.Password)
	fmt.Printf("Password match: %v\n", match)

	if !match {
		return "", errors.New("invalid credentials (password mismatch)")
	}

	return util.GenerateJWT(user.ID)
}
