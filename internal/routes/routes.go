package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
	"transcribify/internal/models"
	"transcribify/internal/routes/middlewares"
	"transcribify/pkg/finders"
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
	if valid, err := middlewares.ValidateVideoRequest(VideoRequest); !valid || err != nil {
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

func (route *Route) HelloWorld(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "Hello World")
}

func (route *Route) SignUp(w http.ResponseWriter, r *http.Request) {

	// Get data from query
	input := route.getSignInData(r)

	// Created user
	err := route.service.Authorization.SignUser(r.Context(), w, input)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		route.logger.Info("Failed to create user", zap.Error(err))

		return
	}
}

func GetSubFromCtx(ctx context.Context) string {
	switch val := ctx.Value("sub"); val {
	case nil:
		return ""
	default:
		return val.(string)
	}
}

func (route *Route) LogIn(w http.ResponseWriter, r *http.Request) {

	// Get data from query
	input := route.getSignInData(r)

	err := route.service.Authorization.LoginUser(r.Context(), w, input)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	route.logger.Info("Set `JWT` token for user",
		zap.Any("user", input))

}

func (route *Route) GetToken(w http.ResponseWriter, r *http.Request) {

}

func (route *Route) getSignInData(r *http.Request) *models.User {
	return &models.User{
		ID:        0,
		Email:     r.URL.Query().Get("login"),
		Password:  r.URL.Query().Get("password"),
		Role:      "",
		CreatedAt: time.Time{},
		LastVisit: time.Time{},
	}

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
