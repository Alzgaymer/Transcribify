package finders

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"yt-video-transcriptor/models"
	"yt-video-transcriptor/models/repository"
)

type (
	// Finder represents interface for founding
	// transcription of the YouTube video.
	Finder interface {
		Find(context.Context, models.VideoRequest) (*models.YTVideo, error)
	}
	APIFinder struct {
		client  *http.Client
		headers http.Header
		repo    repository.Repository
	}
	DatabaseFinder struct {
		repo repository.Repository
	}
)

func NewAPIFinder(client *http.Client, repository repository.Repository) *APIFinder {
	return newAPIFinderWithDefaultHeaders(client, repository)
}

func newAPIFinderWithDefaultHeaders(client *http.Client, repository repository.Repository) *APIFinder {
	headers := http.Header{}

	headers.Add("X-RapidAPI-Key", os.Getenv("VIDEO_API_KEY"))
	headers.Add("X-RapidAPI-Host", os.Getenv("VIDEO_API_URL"))

	return NewAPIFinderWithHeaders(
		client,
		headers,
		repository,
	)
}

func NewAPIFinderWithHeaders(client *http.Client, headers http.Header, repository repository.Repository) *APIFinder {
	return &APIFinder{
		client:  client,
		headers: headers,
		repo:    repository,
	}
}

func NewDatabaseFinder(repository repository.Repository) *DatabaseFinder {
	return &DatabaseFinder{
		repo: repository,
	}
}

func (a *APIFinder) Find(ctx context.Context, video models.VideoRequest) (*models.YTVideo, error) {
	var (
		APIURL = fmt.Sprintf(
			"https://youtube-transcriptor.p.rapidapi.com/transcript?video_id=%s&lang=%s",
			video.VideoID,
			video.Language,
		)
		data []models.YTVideo
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = a.headers

	response, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	err = a.repo.Create(ctx, data[0], video)
	if err != nil {
		return nil, err
	}

	return &data[0], nil
}

func (d *DatabaseFinder) Find(ctx context.Context, video models.VideoRequest) (*models.YTVideo, error) {
	read, err := d.repo.Read(ctx, video)
	if err != nil {
		return nil, err
	}

	return &read, nil
}
