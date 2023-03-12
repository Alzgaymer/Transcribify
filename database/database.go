package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"yt-video-transcriptor/models"
)

// Client has signature of pgx.Tx interface
type Client interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository interface {
	Insert(context.Context, ...models.YTVideo) error
	Read(context.Context, models.YTVideo) ([]byte, error)
	Update(context.Context, models.YTVideo) error
	Delete(context.Context, ...models.YTVideo) error
}
