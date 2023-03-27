package routes

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
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
	)

	//Validating request
	if Valid, err := isValidVideoRequest(VideoRequest); !Valid || err != nil {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}

	for _, finder := range route.finders {
		video, err = finder.Find(r.Context(), VideoRequest)
		if err == nil {
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	writeVideoJson(w, video)
}

// I want you to summarize. I give you a youtube video transcription. You giving me summarizing info, what is going on in the video. Here is transcriptions: %s

func writeVideoJson(w io.Writer, obj any) error {

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	return encoder.Encode(obj)
}

func isValidVideoRequest(request models.VideoRequest) (bool, error) {

	var matchedVideo, err = regexp.MatchString("^[a-zA-Z0-9_-]{11}$", request.VideoID)

	matchedLang, err := regexp.MatchString("^[a-zA-Z]{2}$", request.Language)

	return matchedVideo && matchedLang, err
}
