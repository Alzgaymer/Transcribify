package middlewares

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"log"
	"net/http"
	"transcribify/internal/models"
	"transcribify/pkg/logging"
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
			)

			videoRequest.VideoID = chi.URLParam(r, models.VideoIDTag)
			videoRequest.Language = r.URL.Query().Get(models.LanguageTag)

			logger.Info("Input parameter",
				zap.String("Video ID", videoRequest.VideoID),
				zap.String("Language", videoRequest.Language),
				zap.String("URL", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})

	}
}

func ValidateVideoRequest(request models.VideoRequest) (bool, error) {
	v := validator.New()
	switch err := v.Struct(request); err {
	case nil:
		return true, nil
	default:
		return false, err
	}
}
