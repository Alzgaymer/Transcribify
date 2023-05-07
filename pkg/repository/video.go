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
	var id int

	rawThumb, err := json.Marshal(video.Thumbnails)
	if err != nil {
		return -2, err
	}
	rawTransc, err := json.Marshal(video.Thumbnails)
	if err != nil {
		return -2, err
	}

	_, err = p.client.Exec(ctx, "call put_video($1, $2, $3, $4, $5, $6, $7, $8)",
		video.Title, video.Description, video.AvailableLangs, video.LengthInSeconds,
		rawThumb, rawTransc, request.VideoID, request.Language,
	)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p *YTVideoRepository) GetVideoByIDLang(ctx context.Context, request models.VideoRequest) (*models.YTVideo, error) {
	var (
		rawQuery = `SELECT id, title, description ,available_langs ,length_in_seconds , thumbnails ,transcription 
					FROM video as vd
					WHERE vd.video_id = $1 and
					      vd.language = $2`
		query    = formatQuery(rawQuery)
		rawThumb json.RawMessage
		rawTrans json.RawMessage
		video    models.YTVideo
	)

	err := p.client.QueryRow(ctx, query, request.VideoID, request.Language).
		Scan(&video.Id, &video.Title, &video.Description, &video.AvailableLangs, &video.LengthInSeconds, &rawThumb, &rawTrans)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawThumb, &video.Thumbnails)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawTrans, &video.Transcription)
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
