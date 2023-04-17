package middlewares

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
	"transcribify/internal/models"
	"transcribify/pkg/auth"
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

func CheckCookie(logger *zap.Logger) func(http.Handler) http.Handler {
	var err error

	if logger == nil {
		logger, err = logging.New()
		if err != nil {
			log.Fatal(err)
		}
	}
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Authorization not empty
			if r.Header.Get("Authorization") != "" {
				logger.Info("`Authorization` header provided")
				next.ServeHTTP(w, r)

				return
			}

			token, err := r.Cookie("access")
			if err != nil {
				logger.Info("No `access` JWT token provided in cookie", zap.Error(err))
				next.ServeHTTP(w, r)

				return
			}

			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))

			next.ServeHTTP(w, r)
		})
	}
}

func Identify(logger *zap.Logger, manager auth.TokenManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")

			if header == "" {

				logger.Info("Empty 'Authorization' header",
					zap.String("header", header),
					zap.String("url", r.URL.RawPath))
				next.ServeHTTP(w, r)

				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {

				logger.Info("Invalid 'Authorization' header",
					zap.String("header", header),
					zap.String("url", r.URL.RawPath),
				)
				next.ServeHTTP(w, r)

				return
			}

			if len(headerParts[1]) == 0 {
				logger.Info("Token is empty", zap.String("header", header), zap.String("url", r.URL.RawPath))
				next.ServeHTTP(w, r)

				return
			}

			id, err := manager.Parse(headerParts[1])
			if err != nil {
				logger.Info("Token parsing error", zap.Error(err))
				next.ServeHTTP(w, r)

				return
			}
			logger.Info("Identified", zap.Int("id", id))
			ctx := context.WithValue(r.Context(), "sub", id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
