package database

import (
	"context"
	"yt-video-transcriptor/models"
)

type Repository interface {
	Create(context.Context, models.YTVideo) error
	Read(context.Context, string) (models.YTVideo, error)
	Update(context.Context, string, models.YTVideo) error
	Delete(context.Context, ...models.YTVideo) error
}
