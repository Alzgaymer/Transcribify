package database

import (
	"context"
	"yt-video-transcriptor/models"
)

type Postgres struct {
	client Client
}

func (p Postgres) Insert(ctx context.Context, video ...models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) Read(ctx context.Context, video models.YTVideo) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) Update(ctx context.Context, video models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) Delete(ctx context.Context, video ...models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgres(client Client) Repository {
	return &Postgres{
		client: client,
	}
}
