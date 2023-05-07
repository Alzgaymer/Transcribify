package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
	"transcribify/pkg/hash"
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

		// GetUserByLogin use models.User Email and Password fields to fill model.User struct.
		GetUserByLogin(ctx context.Context, user *models.User) error

		GetUserVideos(ctx context.Context, uid int, limit int, offset int) (map[int]models.YTVideo, error)

		// PutUser store user. If user exist fill models.User ID field.
		PutUser(ctx context.Context, user *models.User) error

		PutUserVideo(ctx context.Context, uid int, vidID int) error
	}
)

func NewRepositories(client *pgx.Conn, hasher hash.PasswordHasher) *Repository {
	return &Repository{
		Video: NewYTVideoRepository(client),
		User:  NewUserRepository(client, hasher),
	}
}
