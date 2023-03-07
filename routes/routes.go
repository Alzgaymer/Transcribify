package routes

import (
	"encoding/json"
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
	videoID, language string
}

// GetVideoTranscription Handle GET request for video with specified language
func (route *Route) GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	// Get the video ID and language from the URL parameters

	videoID := r.URL.Query().Get("v")
	language := r.URL.Query().Get("lang")

	route.logger.Info("Input parameter",
		zap.String("Video ID", videoID),
		zap.String("Language", language),
		zap.String("URL", r.URL.Path),
	)
	videoRequest := VideoRequest{
		videoID:  videoID,
		language: language,
	}

	//Validating request
	if isNotValidVideoRequest(videoRequest) {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}

	//Make request to API
	res, err := route.callbackToApi(videoRequest)
	if err != nil {
		route.logger.Error("Failed to get transcription", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	//Reading response.Body into []model.YTVideo
	video, err := responseToYTVideo(res)
	if err != nil {
		route.logger.Error("Failed to unmarshal data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = writeVideoJson(w, video)
	if err != nil {
		route.logger.Error("Failed to encode data", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (route *Route) callbackToApi(request VideoRequest) (*http.Response, error) {

	configuration := config.GetAPI()

	url := fmt.Sprintf(
		"https://youtube-transcriptor.p.rapidapi.com/transcript?video_id=%s&lang=%s",
		request.videoID,
		request.language,
	)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("X-RapidAPI-Key", configuration.Key)
	req.Header.Add("X-RapidAPI-Host", configuration.API)

	return route.client.Do(req)
}

func responseToYTVideo(res *http.Response) ([]models.YTVideo, error) {

	var data []models.YTVideo
	err := json.NewDecoder(res.Body).Decode(&data)
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

func isValidVideoRequest(request VideoRequest) bool {
	matchedVideo, _ := regexp.MatchString("^[a-zA-Z0-9_-]{11}$", request.videoID)

	matchedLang, _ := regexp.MatchString("^[a-zA-Z]{2}$", request.language)

	return matchedVideo && matchedLang
}

func isNotValidVideoRequest(request VideoRequest) bool {
	return !isValidVideoRequest(request)
}
