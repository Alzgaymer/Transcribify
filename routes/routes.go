package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/models"

	"yt-video-transcriptor/logger"
)

// GetVideoTranscription Handle GET request for video with specified language
func GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	log, err := logger.New(
		logger.WithDevelopment(true),
		logger.WithLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
	)

	if err != nil {
		return
	}

	// Get the video ID and language from the URL parameters
	videoID := chi.URLParam(r, "videoID")
	language := chi.URLParam(r, "language")

	log.Info("Input parameter",
		zap.String("Video ID", videoID),
		zap.String("Language", language),
		zap.String("URL", r.URL.Path),
	)

	configuration := config.GetAPI()

	url := fmt.Sprintf(
		"https://youtube-transcriptor.p.rapidapi.com/transcript?video_id=%s&lang=%s",
		videoID,
		language,
	)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("X-RapidAPI-Key", configuration.Key)
	req.Header.Add("X-RapidAPI-Host", configuration.API)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Error("Failed to get transcription", zap.Error(err))
	}
	defer res.Body.Close()

	log.Info("Getting response...")

	// Use the decoder to parse the response
	var data []models.YTVideo
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Error("Failed to unmarshal data", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	if err != nil {
		log.Error("Failed to encode data", zap.Error(err))
	}
}
