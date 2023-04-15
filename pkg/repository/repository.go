package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
)

const (
	NotFound = -1 * (iota + 1)
	InternalRepositoryError
)

type (
	Repository struct {
		Video Video
		User  User
	}

	Video interface {
		Create(context.Context, models.YTVideo, models.VideoRequest) (int, error)
		Read(context.Context, models.VideoRequest) (models.YTVideo, error)
		Update(context.Context, models.VideoRequest, models.YTVideo) error
		Delete(context.Context, models.VideoRequest) error
	}

	User interface {
		GetUserId(ctx context.Context, login string) (string, error)

		// GetUserByLogin searching in database for user with provided login.
		// If success returns user`s id.
		// If failed - returns "-1", error. "-1" means there`re no users with provided login
		GetUserByLogin(ctx context.Context, login string) (string, error)

		// GetUserByLoginPassword finds user by provided models.User(Email, Password)
		// If user not found - returns NotFound, and user(user.Email) doesn`t exist
		// If internal error - returns InternalRepositoryError, error
		// Else user`s id, nil
		GetUserByLoginPassword(ctx context.Context, user *models.User) (string, error)

		// PutUser firstly search for existing user.
		// If found - returns user`s id, error.
		// If not - inserting. On failure returns InternalRepositoryError, err.
		// On success returns user`s id, nil
		PutUser(ctx context.Context, user *models.User) (string, error)
		SetRefreshToken(ctx context.Context, login, token string) error
		GetRefreshTokenByID(ctx context.Context, id string) (string, error)
	}
)

func convCode(code int) string {
	return fmt.Sprintf("%d", code)
}

func NewRepositories(client *pgx.Conn) *Repository {
	return &Repository{
		Video: NewYTVideoRepository(client),
		User:  NewUserRepository(client),
	}
}
