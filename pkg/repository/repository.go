package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
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
		GetUserId(ctx context.Context, login string) (int, error)

		// SignUser If user with provided login exist returns his id
		// If not - creates in database and returns his id
		SignUser(ctx context.Context, login, password string) (string, error)
		LoginUser(ctx context.Context, login, password string) (string, error)
		SetRefreshToken(ctx context.Context, login, token string) error
		GetRefreshTokenByID(ctx context.Context, id string) (string, error)
	}
)

func NewRepositories(client *pgx.Conn) *Repository {
	return &Repository{
		Video: NewYTVideoRepository(client),
		User:  NewUserRepository(client),
	}
}
