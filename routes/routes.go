package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/database"
	"yt-video-transcriptor/models"
)

type Route struct {
	logger     *zap.Logger
	client     *http.Client
	repository database.Repository
}

func NewRoute(logger *zap.Logger, client *http.Client, repository database.Repository) *Route {
	return &Route{logger: logger, client: client, repository: repository}
}

type VideoRequest struct {
	VideoID  string `json:"v"`
	Language string `json:"lang"`
}

// GetVideoTranscription Handle GET request for video with specified language
func (route *Route) GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	// Get the video ID and language from the query
	videoRequest := VideoRequest{
		VideoID:  r.URL.Query().Get("v"),
		Language: r.URL.Query().Get("lang"),
	}

	//Validating request
	if Valid, err := isValidVideoRequest(videoRequest); !Valid || err != nil {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}

	//Make request to API
	res, err := route.requestToApi(
		"https://youtube-transcriptor.p.rapidapi.com/transcript?video_id=%s&lang=%s",
		videoRequest)
	if err != nil {
		route.logger.Error("Failed to get transcription", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	//Reading response.Body into []model.YTVideo
	video, err := responseToYTVideo(res.Body)
	if err != nil {
		route.logger.Error("Failed to unmarshal data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Writes to html page
	w.WriteHeader(http.StatusOK)
	err = writeVideoJson(w, video)
	if err != nil {
		route.logger.Error("Failed to encode data", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (route *Route) requestToApi(APIurl string, request VideoRequest) (*http.Response, error) {

	configuration := config.API()

	url := fmt.Sprintf(
		APIurl,
		request.VideoID,
		request.Language,
	)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("X-RapidAPI-Key", configuration.Key)
	req.Header.Add("X-RapidAPI-Host", configuration.API)

	return route.client.Do(req)
}

func responseToYTVideo(res io.Reader) ([]models.YTVideo, error) {
	if res == nil {
		return nil, errors.New("io.Reader is nil")
	}

	var data []models.YTVideo
	err := json.NewDecoder(res).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeVideoJson(w io.Writer, video []models.YTVideo) error {

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	return encoder.Encode(video)
}

func isValidVideoRequest(request VideoRequest) (bool, error) {
	matchedVideo, err := regexp.MatchString("^[a-zA-Z0-9_-]{11}$", request.VideoID)

	matchedLang, err := regexp.MatchString("^[a-zA-Z]{2}$", request.Language)

	return matchedVideo && matchedLang, err
}
