package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yt-video-transcriptor/internal/config"
	"yt-video-transcriptor/internal/logger"
)

func main() {
	// Setting up server configuration (port etc.)
	var configuration config.AppConfiguration
	err := cleanenv.ReadConfig("route.env", &configuration)
	if err != nil {
		panic(err)
	}

	server := http.Server{Addr: ":" + configuration.Port, Handler: service()}

	log, err := logger.New(
		logger.WithDevelopment(true),
		logger.WithLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
	)

	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Waiting to signal from os
	go func() {
		defer cancel()
		<-sig

		timeout, _ := context.WithTimeout(ctx, 30*time.Second)

		go func() {
			// Waiting until context is done then panic
			<-timeout.Done()
			if timeout.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out... forcing exit...")
			}
		}()

		log.Info("Stopping server")

		err = server.Shutdown(timeout)
		if err != nil {
			log.Fatal("Failed to shutdown server", zap.Error(err))
		}

	}()

	log.Info("Server is running", zap.String("port", configuration.Port))

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to serve", zap.Error(err))
	}

	<-ctx.Done()
}

func service() http.Handler {

	log, err := logger.New()
	if err != nil {
		return nil
	}

	router := chi.NewRouter()

	router.Get("/api/v1/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("<h1>Hello</h1>"))
	})

	videoIDPattern := "[a-zA-Z0-9_-]{11}"

	// Create a route for the GET method that accepts the video ID as a parameter
	router.Get("/api/v1/{videoID:"+videoIDPattern+"}", func(w http.ResponseWriter, r *http.Request) {

		// Get the video ID from the URL parameters
		videoID := chi.URLParam(r, "videoID")

		// Do something with the video ID, for example print it to the console
		log.Info("Input parametr",
			zap.String("Video ID", videoID),
			zap.String("URL", "/api/v1/"+videoID),
		)

		// Return a response to the client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Video ID: " + videoID))
	})
	return router
}
