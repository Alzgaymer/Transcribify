package auth

import (
	"context"
	"transcribify/internal/models"
	"transcribify/pkg/repository"
)

type Authorization interface {
	// SignUser If user with provided login exist returns his id
	// If not - creates in database and returns his id
	SignUser(ctx context.Context, user *models.User) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

type AuthorizationManager struct {
	repository repository.User
}

func (a *AuthorizationManager) SignUser(ctx context.Context, user *models.User) (string, error) {
	return a.repository.PutUser(ctx, user)
}

func (a *AuthorizationManager) LoginUser(ctx context.Context, login, password string) (string, error) {
	return a.repository.GetUserByLoginPassword(ctx, &models.User{Password: password, Email: login})
}

func NewAuthorizationManager(user repository.User) *AuthorizationManager {
	return &AuthorizationManager{repository: user}
}
