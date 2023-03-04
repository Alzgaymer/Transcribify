package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/models"

	"yt-video-transcriptor/logger"
)

// GetVideoTranscription Handle GET request for video with specified language
func GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	log, err := logger.New(
		logger.WithDevelopment(true),
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

	log.Info("Getting response...")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Failed to get transcription", zap.Error(err))
	}
	defer res.Body.Close()

	//Reading response.Body into []model.YTVideo
	video, err := responseToYTVideo(res)
	if err != nil {
		log.Error("Failed to unmarshal data", zap.Error(err))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = writeVideoJson(w, video)
	if err != nil {
		log.Error("Failed to encode data", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)

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
