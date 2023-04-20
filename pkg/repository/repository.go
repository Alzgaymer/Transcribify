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
		CreateVideo(context.Context, models.VideoRequest, *models.YTVideo) (int, error)

		GetVideoByIDLang(context.Context, models.VideoRequest) (*models.YTVideo, error)

		Update(context.Context, models.VideoRequest, *models.YTVideo) error
		Remove(context.Context, models.VideoRequest) error
	}

	User interface {

		// GetUserByLoginPassword use models.User Email and Password fields to fill model.User struct.
		GetUserByLoginPassword(ctx context.Context, user *models.User) error

		GetUserVideos(ctx context.Context, uid int, vid string) ([]string, error)

		// PutUser store user. If user exist fill models.User ID field.
		PutUser(ctx context.Context, user *models.User) error

		PutUserVideo(ctx context.Context, uid int, vid string) error
	}
)

func NewRepositories(client *pgx.Conn) *Repository {
	return &Repository{
		Video: NewYTVideoRepository(client),
		User:  NewUserRepository(client),
	}
}
