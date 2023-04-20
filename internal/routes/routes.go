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
	"transcribify/pkg/repository"
	"transcribify/pkg/service"
)

type Route struct {
	logger     *zap.Logger
	client     *http.Client
	repository *repository.Repository
	service    *service.Services
}

func NewRoute(
	logger *zap.Logger,
	client *http.Client,
	repository *repository.Repository,
	service *service.Services,
) *Route {
	return &Route{
		logger:     logger,
		client:     client,
		repository: repository,
		service:    service,
	}

}

// GetVideoTranscription Handle GET request for video with specified language
func (route *Route) GetVideoTranscription(w http.ResponseWriter, r *http.Request) {

	var (
		vr = models.VideoRequest{
			VideoID:  chi.URLParam(r, "id"),
			Language: r.URL.Query().Get("lang"),
		}
		video = new(models.YTVideo)
		err   error
		ctx   = r.Context()
	)
	// Get the language from the query
	uid := GetSubFromCtx(ctx)

	if uid == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		route.logger.Info("Invalid user id", zap.Int("uid", uid))

		return
	}

	//Validating request
	if valid, err := middlewares.ValidateVideoRequest(vr); !valid || err != nil {
		w.WriteHeader(http.StatusConflict)
		route.logger.Info("Invalid video request",
			zap.Any("video request", vr), zap.Error(err), zap.Bool("valid", valid))
		return
	}

	video, err = route.service.Finder.Find(ctx, vr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to find video", zap.Error(err))

		return
	}

	err = route.repository.User.PutUserVideo(ctx, uid, vr.VideoID)
	if err != nil {
		route.logger.Info("Failed to put user video", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

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

// GetSubFromCtx returns -1 if sub doesn`t provided in context.Context
func GetSubFromCtx(ctx context.Context) int {
	switch val := ctx.Value("sub"); val {
	case nil:
		return -1
	default:
		return val.(int)
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

func (route *Route) GetUserVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user id
	uid := GetSubFromCtx(ctx)

	if uid == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		route.logger.Info("Failed to authorize", zap.Int("uid", uid))
		return
	}

	vid := chi.URLParam(r, "vid")

	videos, err := route.repository.User.GetUserVideos(ctx, uid, vid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to get user videos", zap.Error(err))

		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, videos)
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
