package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/database"
	"yt-video-transcriptor/logging"
	"yt-video-transcriptor/routes"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	var (
		configuration = config.GetRoute()
		client        = &http.Client{
			Timeout: 30 * time.Second,
		}
		ctx, cancel                     = context.WithCancel(context.Background())
		logger                          = withLogger()
		postgres    database.Repository = nil //withDatabase(ctx, 5, 1*time.Second)
	)

	server := http.Server{Addr: ":" + configuration.Port, Handler: service(
		logger,
		client,
		postgres,
	)}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Waiting to signal from os
	go func() {
		defer cancel()
		<-sig

		timeout, timeoutcancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutcancel()
		go func() {
			// Waiting until context is done then panic
			<-timeout.Done()
			if errors.Is(timeout.Err(), context.DeadlineExceeded) {
				logger.Fatal("Graceful shutdown timed out... forcing exit...")
			}
		}()

		logger.Info("Stopping server")

		err = server.Shutdown(timeout)
		if err != nil {
			logger.Fatal("Failed to shutdown server", zap.Error(err))
		}
	}()

	logger.Info("Server is running", zap.String("port", configuration.Port))

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to serve", zap.Error(err))
	}

	<-ctx.Done()
}

func service(logger *zap.Logger, client *http.Client, repository database.Repository) http.Handler {

	router := chi.NewRouter()

	//router.Use(logging.Logging(logger))

	route := routes.NewRoute(logger, client, repository)
	// Create a route for the GET method that accepts the video ID as a parameter
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/video", route.GetVideoTranscription) //GET	/api/v1/video?v=&lang=
	})

	return router
}

func withLogger() *zap.Logger {
	logger, err := logging.New(
		logging.WithDevelopment(true),
		logging.WithLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	return logger
}
func withDatabase(ctx context.Context, attemptsToConnect uint, sleep time.Duration) database.Repository {
	client, err := database.NewClient(ctx, attemptsToConnect, sleep)
	if err != nil {
		log.Fatal(err)
	}
	return database.NewPostgres(client)
}
