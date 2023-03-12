package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"yt-video-transcriptor/models"
)

type Postgres struct {
	client *pgx.Conn
}

func NewPostgres(client *pgx.Conn) Repository {
	return &Postgres{
		client: client,
	}
}

func (p Postgres) Create(ctx context.Context, video models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) Read(ctx context.Context, s string) (models.YTVideo, error) {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) Update(ctx context.Context, s string, video models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) Delete(ctx context.Context, video ...models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}
