package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
	"transcribify/pkg/hash"
)

type UserRepository struct {
	client *pgx.Conn
	hash   hash.PasswordHasher
}

func (u *UserRepository) PutUser(ctx context.Context, user *models.User) error {
	return u.client.QueryRow(ctx, "select id from put_user($1, $2)", user.Email, user.Password).
		Scan(&user.ID)
}

func (u *UserRepository) GetUserVideos(ctx context.Context, uid int, vid string) ([]string, error) {

	arr := make([]string, 0)
	rows, err := u.client.Query(ctx, "select video_id from get_user_videos($1, $2)", uid, vid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var vid string
		err = rows.Scan(&vid)
		if err != nil {
			return nil, err
		}

		arr = append(arr, vid)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return arr, nil
}

func (u *UserRepository) PutUserVideo(ctx context.Context, uid int, vid string) error {
	_, err := u.client.Exec(ctx, "call put_user_video($1, $2)", uid, vid)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetUserByLoginPassword(ctx context.Context, user *models.User) error {
	return u.client.QueryRow(ctx, "select id, password from get_user($1, $2)", user.Email, user.Password).
		Scan(&user.ID, &user.Password)
}

func NewUserRepository(client *pgx.Conn) *UserRepository {
	return &UserRepository{client: client}
}
