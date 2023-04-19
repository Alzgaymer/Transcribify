package repository

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
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

func (p *YTVideoRepository) CreateVideo(ctx context.Context, request models.VideoRequest, video *models.YTVideo) (int, error) {
	var (
		id        int
		sb        = new(strings.Builder)
		rawThumb  json.RawMessage
		rawTransc json.RawMessage
	)

	err := json.NewEncoder(sb).Encode(rawThumb)
	if err != nil {
		return -2, err
	}

	sb.Reset()
	err = json.NewEncoder(sb).Encode(rawThumb)
	if err != nil {
		return -2, err
	}

	err = p.client.QueryRow(ctx, "select id from put_video($1, $2, $3, $4, $5, $6, $7, $8)",
		video.Title, video.Description, video.AvailableLangs, video.LengthInSeconds,
		rawThumb, rawTransc, request.VideoID, request.Language,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p *YTVideoRepository) GetVideoByIDLang(ctx context.Context, request models.VideoRequest) (*models.YTVideo, error) {
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
		return nil, err
	}

	err = json.Unmarshal(rawVideo, &video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (p *YTVideoRepository) Update(ctx context.Context, req models.VideoRequest, video *models.YTVideo) error {
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

func (p *YTVideoRepository) Remove(ctx context.Context, req models.VideoRequest) error {
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
