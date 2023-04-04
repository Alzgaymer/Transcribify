package middlewares

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"regexp"
	"yt-video-transcriptor/logging"
	"yt-video-transcriptor/models"
)

func LogVideoRequest(logger *zap.Logger) func(next http.Handler) http.Handler {
	var err error

	if logger == nil {
		logger, err = logging.New()
		if err != nil {
			log.Fatal(err)
		}
	}

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				videoRequest = models.VideoRequest{}
				str          = "Input parameter"
			)

			videoRequest.VideoID = chi.URLParam(r, models.VideoIDTag)
			videoRequest.Language = r.URL.Query().Get(models.LanguageTag)

			if valid, _ := ValidateVideoRequest(videoRequest); !valid {
				str = "Invalid video request"
			}

			logger.Info(str,
				zap.String("Video ID", videoRequest.VideoID),
				zap.String("Language", videoRequest.Language),
				zap.String("URL", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})

	}
}
func ValidateVideoRequest(request models.VideoRequest) (bool, error) {

	var matchedVideo, err = regexp.MatchString(models.VideoPattern, request.VideoID)

	matchedLang, err := regexp.MatchString(models.LanguagePattern, request.Language)

	return matchedVideo && matchedLang, err
}
