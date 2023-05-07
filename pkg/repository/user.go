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
	p := u.hash.Hash(user.Password)
	_, err := u.client.Exec(ctx, "call put_user($1, $2)", user.Email, p)
	return err

}

func (u *UserRepository) GetUserVideos(ctx context.Context, uid int, limit int, offset int) (map[int]models.YTVideo, error) {

	arr := make(map[int]models.YTVideo, 0)
	rows, err := u.client.Query(ctx, "select uv.id, v.title, v.length_in_seconds from public.video v left join user_videos uv on v.id = uv.video_id and uv.user_id = $1 limit $2 offset $3", uid, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var video models.YTVideo
		err = rows.Scan(&id, &video.Title, &video.LengthInSeconds)

		if err != nil {
			return nil, err
		}

		arr[id] = video
	}

	if rows.Err() != nil {
		return nil, err
	}

	return arr, nil
}

func (u *UserRepository) PutUserVideo(ctx context.Context, uid int, vidID int) error {
	_, err := u.client.Exec(ctx, "insert into user_videos(user_id, video_id) values ($1, $2)", uid, vidID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, user *models.User) error {
	return u.client.QueryRow(ctx, "select id, password from get_user($1)", user.Email).
		Scan(&user.ID, &user.Password)
}

func NewUserRepository(client *pgx.Conn, haser hash.PasswordHasher) *UserRepository {
	return &UserRepository{client: client, hash: haser}
}
