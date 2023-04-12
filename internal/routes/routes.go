package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"time"
	"transcribify/internal/models"
	"transcribify/internal/routes/middlewares"
	"transcribify/pkg/auth"
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

			return
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

		id, err := route.service.Manager.Parse(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			route.logger.Info("Token parsing error", zap.Error(err))

			return
		}

		ctx := context.WithValue(r.Context(), "sub", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (route *Route) HelloWorld(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "Hello World")
}

func (route *Route) CheckCookie(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Authorization not empty
		if r.Header.Get("Authorization") != "" {
			route.logger.Info("`Authorization` header provided")
			next.ServeHTTP(w, r)

			return
		}

		token, err := r.Cookie("access")
		if err != nil {
			route.logger.Info("No `access` JWT token provided in cookie", zap.Error(err))
			next.ServeHTTP(w, r)

			return
		}

		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))

		next.ServeHTTP(w, r)
	})
}

func (route *Route) SignUp(w http.ResponseWriter, r *http.Request) {

	// Check Authorization header
	// If exist parse
	if header := r.Header.Get("Authorization"); header != "" {

		route.logger.Info("", zap.String("header", header))

		headerParts := strings.Split(header, " ")
		var toParse string
		if len(headerParts) == 2 {
			toParse = headerParts[1]
		} else {
			toParse = header
		}
		id, err := route.service.Manager.Parse(toParse)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			route.logger.Info("Failed to parse token", zap.String("header", header))

			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("already authorized as user: %s", id)))
		if err != nil {
			route.logger.Info("Failed to write into html", zap.Error(err))
		}

		return
	}

	// If not - repository.SignIn
	// Get data from query
	input := route.getSignInData(r)

	// Created user
	user, err := route.repository.User.SignUser(r.Context(),
		input.Email,
		hash.NewSHA1PSHasher(os.Getenv("JWT_SALT")).Hash(input.Password))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		route.logger.Info("Failed to create user", zap.Error(err))

		return
	}
	route.logger.Info("Get user id", zap.String("user", user))

	//Create token
	jwt, err := route.service.Manager.NewJWT(&models.User{ID: user}, auth.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to create `JWT` token for user", zap.String("user", user))

		return
	}

	SetCookie(w, "access", jwt)
	route.logger.Info("Set `JWT` token for user", zap.String("user", user), zap.String("jwt", jwt.T))
}

func (route *Route) GetToken(w http.ResponseWriter, r *http.Request) {

	input := route.getSignInData(r)

	user, err := route.repository.User.GetUserId(r.Context(), input.Email)

	jwt, err := route.service.Manager.NewJWT(input, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to generate JWT token", zap.Int("user", user), zap.Error(err))

		return
	}
	w.WriteHeader(http.StatusCreated)
	SetCookie(w, "access", jwt)

}

func (route *Route) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("sub").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to get token from context")
		return
	}

	token, err := route.repository.User.GetRefreshTokenByID(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to get refresh token from database", zap.Error(err))

		return
	}

	id, err := route.service.Manager.Parse(token)
	if userId != id {
		w.WriteHeader(http.StatusUnauthorized)
		route.logger.Info("Failed to parse refresh token", zap.Error(err))

		return
	}

	newJWT, err := route.service.Manager.NewJWT(&models.User{Role: "admin", ID: userId}, auth.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		route.logger.Info("Failed to generate new JWT token", zap.Error(err))

		return
	}

	render.JSON(w, r, newJWT)
}

func (route *Route) getSignInData(r *http.Request) *models.User {
	return &models.User{
		ID:        "",
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

func SetCookie(w http.ResponseWriter, name string, token models.Token) {

	cookie := &http.Cookie{
		Name:     name,
		Value:    token.T,
		Expires:  token.ExpiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
