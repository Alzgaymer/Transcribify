package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"transcribify/internal/models"
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks yt-video-transcriptor/models/repository Repository
type (
	YTVideoRepository struct {
		client *pgx.Conn
	}
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
		return InternalRepositoryError, err
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

			return InternalRepositoryError, newErr
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

func (p *YTVideoRepository) Update(ctx context.Context, req models.VideoRequest, video models.YTVideo) error {
	// write a query to update the video

	var (
		rawQuery = `	UPDATE video_data
						SET json_data = $1
						WHERE video_id = $2 AND language = $3;
					`
		query    = formatQuery(rawQuery)
		rawVideo strings.Builder
	)
	//Encode video in json
	if err := json.NewEncoder(&rawVideo).Encode(video); err != nil {
		return err
	}

	row := p.client.QueryRow(ctx, query, rawVideo.String(), req.VideoID, req.Language)

	if err := row.Scan(nil); err != pgx.ErrNoRows {
		return err
	}

	return nil
}

func (p *YTVideoRepository) Delete(ctx context.Context, req models.VideoRequest) error {
	var (
		rawQuerry = `	DELETE FROM video_data
						WHERE video_id = $1 AND language = $2;
					`
		query = formatQuery(rawQuerry)
	)

	if err := p.client.QueryRow(ctx, query, req.VideoID, req.Language).Scan(nil); err != pgx.ErrNoRows {
		return err
	}

	return nil
}
