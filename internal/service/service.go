package service

import (
	"database/sql"
	"errors"

	"restapi/internal/model"
	"restapi/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(r *repository.Repository) *AuthService {
	return &AuthService{repo: r}
}
func (s *AuthService) SignUp(email, password, name string) (*model.User, error) {
	user, err := s.repo.CreateUser(email, password, name)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *AuthService) Login(email, password string) (*model.User, string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}

	token := "fake-jwt-token1"

	return user, token, nil
}
