package routes

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"yt-video-transcriptor/finders"
	"yt-video-transcriptor/models"
	"yt-video-transcriptor/models/repository"
)

type Route struct {
	logger     *zap.Logger
	client     *http.Client
	repository repository.Repository
	finders    []finders.Finder
}

func NewRoute(
	logger *zap.Logger,
	client *http.Client,
	repository repository.Repository,
	finders ...finders.Finder) *Route {
	return &Route{logger: logger, client: client, repository: repository, finders: finders}
}

// GetVideoTranscription Handle GET request for video with specified language
func (route *Route) GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	// Get the video ID and language from the query

	var (
		VideoRequest = models.VideoRequest{
			VideoID:  r.URL.Query().Get("v"),
			Language: r.URL.Query().Get("lang"),
		}
		video *models.YTVideo
		err   error
		ctx   = r.Context()
	)

	//Validating request
	if Valid, err := isValidVideoRequest(VideoRequest); !Valid || err != nil {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}

	for _, finder := range route.finders {
		video, err = finder.Find(ctx, VideoRequest)
		if err == nil && video != nil {
			break
		}
	}
	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	prompt, err := formatPrompt(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		route.logger.Info("Error while formatting video to OPENAI prompt", zap.Error(err))
		return
	}
	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 2,
		Prompt:    prompt,
	}
	resp, err := c.CreateCompletion(ctx, req) // HTTP 400 model`s max tokens 2048 in prompt  ~11`000
	if err != nil {
		route.logger.Info("Error while sending request to OPENAI", zap.Error(err))
		return
	}
	route.logger.Info("OPENAI response", zap.Any("resp", resp))
	w.WriteHeader(http.StatusOK)
	writeJson(w, resp.Choices[0].Text)
}

func writeJson(w io.Writer, obj any) error {

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	return encoder.Encode(obj)
}

func isValidVideoRequest(request models.VideoRequest) (bool, error) {

	var matchedVideo, err = regexp.MatchString("^[a-zA-Z0-9_-]{11}$", request.VideoID)

	matchedLang, err := regexp.MatchString("^[a-zA-Z]{2}$", request.Language)

	return matchedVideo && matchedLang, err
}

func formatPrompt(video *models.YTVideo) (string, error) {
	var (
		promt    = "I want you to summarize. I give you a youtube video transcription. You giving me summarizing info, what is going on in the video. Here is transcriptions: %s"
		toInsert strings.Builder
	)
	err := json.NewEncoder(&toInsert).Encode(video)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(promt, toInsert.String()), nil
}
