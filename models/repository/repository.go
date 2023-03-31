package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"yt-video-transcriptor/models"
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks yt-video-transcriptor/models/repository Repository
type (
	Repository interface {
		Create(context.Context, models.YTVideo, models.VideoRequest) (int, error)
		Read(context.Context, models.VideoRequest) (models.YTVideo, error)
		Update(context.Context, string, models.YTVideo) error
		Delete(context.Context, string) error
	}
	YTVideoRepository struct {
		client *pgx.Conn
	}
)

const (
	NotFound                = -1
	RepositoryInternalError = -2
)

func NewYTVideoRepository(client *pgx.Conn) *YTVideoRepository {
	return &YTVideoRepository{
		client: client,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (p *YTVideoRepository) Create(ctx context.Context, video models.YTVideo, request models.VideoRequest) (int, error) {

	var (
		rawQuery = `INSERT INTO video_data(
	                       video_id,
	                       language,
	                       json_data
	                       )
					VALUES ($1, $2, $3)
					returning id`
		query = formatQuery(rawQuery)
		id    int
	)

	jsonData, err := json.Marshal(video)
	if err != nil {
		return RepositoryInternalError, err
	}

	err = p.client.QueryRow(ctx,
		query,
		request.VideoID,
		request.Language,
		jsonData,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))

			return RepositoryInternalError, newErr
		}

		return NotFound, err
	}

	return id, nil
}

func (p *YTVideoRepository) Read(ctx context.Context, request models.VideoRequest) (models.YTVideo, error) {
	var (
		rawQuery = `SELECT json_data 
					FROM video_data as vd
					WHERE vd.video_id = $1 and
					      vd.language = $2`
		query    = formatQuery(rawQuery)
		rawVideo json.RawMessage
		video    models.YTVideo
	)

	row := p.client.QueryRow(ctx, query, request.VideoID, request.Language)

	err := row.Scan(&rawVideo)
	if err != nil {
		return video, err
	}

	err = json.Unmarshal(rawVideo, &video)
	if err != nil {
		return video, err
	}

	return video, nil
}

func (p *YTVideoRepository) Update(ctx context.Context, s string, video models.YTVideo) error {
	//TODO implement me
	panic("implement me")
}

func (p *YTVideoRepository) Delete(ctx context.Context, video string) error {
	//TODO implement me
	panic("implement me")
}
