package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
	"transcribify/internal/models"
	"transcribify/internal/routes"
	"transcribify/internal/routes/middlewares"
	"transcribify/pkg/dbclient"
	"transcribify/pkg/logging"
	repo "transcribify/pkg/repository"
	"transcribify/pkg/service"
)

func Server(ctx context.Context) *http.Server {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return &http.Server{
		Addr:    ":" + os.Getenv("APP_PORT"),
		Handler: Router(Logger(), Client(), service.New(), Repository(ctx)),
	}
}

func Repository(ctx context.Context) *repo.Repository {

	client, err := dbclient.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return repo.NewRepositories(client)
}

// Router uses for http.Server struct Handler field.
// It implements several endpoints
func Router(logger *zap.Logger, client *http.Client, service *service.Service, repository *repo.Repository) http.Handler {

	router := chi.NewRouter()

	route := routes.NewRoute(
		logger, client, repository, service, dbclient.CacheVideoFinders(client, repository.Video)...,
	)

	// Create a route for the GET method that accepts the video ID as a parameter
	router.Route("/api/v1", func(r chi.Router) {
		//GET	/api/v1/{videoID:^[a-zA-Z0-9_-]{11}$}?lang=
		r.With(middlewares.LogVideoRequest(logger)).
			Get(fmt.Sprintf("/{%s:%s}", models.VideoIDTag, models.VideoPattern), route.GetVideoTranscription)

		r.With(route.IdentifyUser).Group(func(r chi.Router) {
			r.Get("/token", route.GetToken)
			r.Get("/refresh", route.RefreshToken)
		})
	})

	return router
}

func Logger() *zap.Logger {

	logg, err := logging.New(
		logging.WithDevelopment(true),
		logging.WithLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return logg
}

// Client returns pointer to the http.Client with 30s timeout
func Client() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
