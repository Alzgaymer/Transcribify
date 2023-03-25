package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/models"
	"yt-video-transcriptor/models/repository"
)

type Route struct {
	logger     *zap.Logger
	client     *http.Client
	repository repository.Repository
}

func NewRoute(logger *zap.Logger, client *http.Client, repository repository.Repository) *Route {
	return &Route{logger: logger, client: client, repository: repository}
}

// GetVideoTranscription Handle GET request for video with specified language
func (route *Route) GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	// Get the video ID and language from the query
	video := models.YTVideo{
		VideoRequest: models.VideoRequest{
			VideoID:  r.URL.Query().Get("v"),
			Language: r.URL.Query().Get("lang"),
		},
	}

	//Validating request
	if Valid, err := isValidVideoRequest(video.VideoRequest); !Valid || err != nil {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}

	//Find in repository
	read, err := route.repository.Read(r.Context(), video.VideoRequest)
	if err == nil && len(read.Transcription) != 0 {
		w.WriteHeader(http.StatusOK)

		err = writeVideoJson(w, read)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
	route.logger.Info("Failed to find", zap.Error(err))

	//Make request to API
	res, err := route.requestToApi(
		http.MethodGet, fmt.Sprintf(
			"https://youtube-transcriptor.p.rapidapi.com/transcript?video_id=%s&lang=%s",
			video.VideoRequest.VideoID,
			video.VideoRequest.Language,
		),
		nil,
		[]config.APIConfiguration{
			{
				Header: "X-RapidAPI-Value", Value: os.Getenv("VIDEO_API_KEY"),
			},
			{
				Header: "X-RapidAPI-Host", Value: os.Getenv("VIDEO_API_URL"),
			},
		})
	if err != nil {
		route.logger.Info("Failed to get transcription", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	//Reading response.Body into model.YTVideo
	tempVideo, err := responseToYTVideo(res.Body)
	if err != nil {
		route.logger.Info("Failed to unmarshal data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tempVideo.VideoRequest = video.VideoRequest
	video = tempVideo

	err = route.repository.Create(r.Context(), video)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		route.logger.Info("Failed to create video_data instance", zap.Error(err))
		return
	}

	//request to openai
	body, err := formatBody(&video)
	if err != nil {
		return
	}
	openaiRes, err := route.requestToApi(
		http.MethodPost,
		"https://api.openai.com/v1/completions",
		body,
		[]config.APIConfiguration{
			{
				Header: "Authorization", Value: "Bearer " + os.Getenv("OPENAI_API_KEY"),
			},
			//OpenAI-Organization org-P8QTzGDPgjPTZ9CCMwB1Qzq8
		})
	if err != nil || openaiRes.Status != "200 OK" {

		route.logger.Info("Failed to request to openai", zap.Error(err),
			zap.String("Status", openaiRes.Status))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer openaiRes.Body.Close()

	w.WriteHeader(http.StatusOK)
	var out []byte
	_, err = openaiRes.Body.Read(out)
	if err != nil {
		route.logger.Info("Failed to read response body", zap.Error(err))
		return
	}
	err = writeVideoJson(w, out)
	if err != nil {
		route.logger.Error("Failed to write to html", zap.Error(err))
		return
	}
}
func formatBody(video *models.YTVideo) (io.Reader, error) {
	str := "{\n\t\t\"model\": \"text-davinci-003\",\n\t\t\"prompt\": \"I want you to summarize. I give you a youtube video transcription. You giving me summarizing info, what is going on in the video. Here is transcriptions: %s\",\n\t\t\"temperature\": 0.1,\n\t\t\"max_tokens\": 512,\n\t\t\"top_p\": 1,\n\t\t\"frequency_penalty\": 0,\n\t\t\"presence_penalty\": 0\n\t}"
	var sb strings.Builder

	for _, transcription := range video.Transcription {
		err := json.NewEncoder(&sb).Encode(transcription)
		if err != nil {
			return nil, err
		}
	}
	return strings.NewReader(fmt.Sprintf(str, sb.String())), nil
}
func (route *Route) requestToApi(method string, uri string, body io.Reader, configurations []config.APIConfiguration) (*http.Response, error) {

	req, _ := http.NewRequest(method, uri, body)

	for _, configuration := range configurations {
		req.Header.Add(configuration.Header, configuration.Value)
	}

	return route.client.Do(req)
}

func responseToYTVideo(res io.Reader) (models.YTVideo, error) {
	if res == nil {
		return models.YTVideo{}, errors.New("io.Reader is nil")
	}

	var data []models.YTVideo
	err := json.NewDecoder(res).Decode(&data)
	if err != nil {
		return models.YTVideo{}, err
	}

	return data[0], nil
}

func writeVideoJson(w io.Writer, obj any) error {

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	return encoder.Encode(obj)
}

func isValidVideoRequest(request models.VideoRequest) (bool, error) {
	matchedVideo, err := regexp.MatchString("^[a-zA-Z0-9_-]{11}$", request.VideoID)

	matchedLang, err := regexp.MatchString("^[a-zA-Z]{2}$", request.Language)

	return matchedVideo && matchedLang, err
}
