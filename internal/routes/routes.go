package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"transcribify/internal/models"
	"transcribify/internal/routes/middlewares"
	"transcribify/pkg/finders"
	"transcribify/pkg/hash"
	"transcribify/pkg/repository"
	"transcribify/pkg/service"
)

type Route struct {
	logger     *zap.Logger
	client     *http.Client
	repository *repository.Repository
	finders    []finders.Finder
	service    *service.Service
}

func NewRoute(
	logger *zap.Logger,
	client *http.Client,
	repository *repository.Repository,
	service *service.Service,
	finders ...finders.Finder,
) *Route {
	return &Route{
		logger:     logger,
		client:     client,
		repository: repository,
		service:    service,
		finders:    finders,
	}
}

// GetVideoTranscription Handle GET request for video with specified language
func (route *Route) GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	// Get the language from the query

	var (
		VideoRequest = models.VideoRequest{
			VideoID:  chi.URLParam(r, models.VideoIDTag),
			Language: r.URL.Query().Get("lang"),
		}
		video *models.YTVideo
		err   error
		ctx   = r.Context()
	)

	//Validating request
	if Valid, err := middlewares.ValidateVideoRequest(VideoRequest); !Valid || err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	for _, finder := range route.finders {
		video, err = finder.Find(ctx, VideoRequest)
		if err == nil && video != nil {
			break
		}
	}

	//c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	//prompt, err := formatPrompt(video)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write(nil)
	//	route.logger.Info("Error while formatting video to OPENAI prompt", zap.Error(err))
	//	return
	//}
	//req := openai.CompletionRequest{
	//	Model:     openai.GPT3Ada,
	//	MaxTokens: 2,
	//	Prompt:    prompt,
	//}
	//resp, err := c.CreateCompletion(ctx, req) // HTTP 400 model`s max tokens 2048 in prompt  ~11`000
	//if err != nil {
	//	route.logger.Info("Error while sending request to OPENAI", zap.Error(err))
	//	return
	//}
	//route.logger.Info("OPENAI response", zap.Any("resp", resp))
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, video)
}

// IdentifyUser is a middleware for identifying user token
func (route *Route) IdentifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")

		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			route.logger.Info("Empty 'Authorization' header", zap.String("header", header), zap.String("url", r.URL.RawPath))

		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			route.logger.Info("Invalid 'Authorization' header", zap.String("header", header), zap.String("url", r.URL.RawPath))
			return
		}

		if len(headerParts[1]) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			route.logger.Info("Token is empty", zap.String("header", header), zap.String("url", r.URL.RawPath))
			return
		}

		parse, err := route.service.Manager.Parse(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			route.logger.Info("Token parsing error", zap.Error(err))

			return
		}

		render.JSON(w, r, parse)

		next.ServeHTTP(w, r)
	})
}

func (route *Route) GetToken(w http.ResponseWriter, r *http.Request) {
	pwhasher := hash.NewSHA1PSHasher(os.Getenv("JWT_SALT"))
	var input models.SignIn

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to unmarshal request `Body`", zap.Error(err))

		return
	}
	pwhasher.Hash(input.Password)
}

func (route *Route) RefreshToken(w http.ResponseWriter, r *http.Request) {

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
