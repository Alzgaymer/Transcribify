package auth

import (
	"context"
	"net/http"
	"transcribify/internal/models"
	"transcribify/pkg/repository"
)

type Authorization interface {
	//Authorize(r *http.Request) error
	SignUser(ctx context.Context, w http.ResponseWriter, user *models.User) error
	LoginUser(ctx context.Context, w http.ResponseWriter, user *models.User) error
}

type AuthorizationManager struct {
	repository repository.User
	tm         TokenManager
}

func NewAuthorizationManager(repository repository.User, tm TokenManager) *AuthorizationManager {
	return &AuthorizationManager{repository: repository, tm: tm}
}

//func (a *AuthorizationManager) Authorize(r *http.Request) error {
//	re
//}

func (a *AuthorizationManager) SignUser(ctx context.Context, w http.ResponseWriter, user *models.User) error {
	err := a.repository.PutUser(ctx, user)
	if err != nil {
		return err
	}

	access, err := a.tm.NewJWT(user, Access)
	if err != nil {
		return err
	}

	refresh, err := a.tm.NewJWT(user, Refresh)
	if err != nil {
		return err
	}

	SetJwtToCookie(w, access, refresh)

	return nil

}

func (a *AuthorizationManager) LoginUser(ctx context.Context, w http.ResponseWriter, user *models.User) error {
	err := a.repository.GetUserByLoginPassword(ctx, user)
	if err != nil {
		return err
	}

	access, err := a.tm.NewJWT(user, Access)
	if err != nil {
		return err
	}

	refresh, err := a.tm.NewJWT(user, Refresh)
	if err != nil {
		return err
	}

	SetJwtToCookie(w, access, refresh)

	return nil
}
func SetJwtToCookie(w http.ResponseWriter, tokens ...models.Token) {

	for _, token := range tokens {
		cookie := &http.Cookie{
			Name:     token.Key,
			Value:    token.T,
			Expires:  token.ExpiresAt,
			HttpOnly: true,
			Path:     "/api/v1/",
		}
		http.SetCookie(w, cookie)
	}
}
