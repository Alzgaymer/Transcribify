package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/logger"
	"yt-video-transcriptor/routes"
)

func main() {
	// Setting up server configuration (port etc.)
	configuration := config.GetRoute()

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

		timeout, timeoutcancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutcancel()
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

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	// Create a route for the GET method that accepts the video ID as a parameter
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/{videoID:[a-zA-Z0-9_-]{11}}", func(r chi.Router) {
			r.Get("/{language:[a-zA-Z]{2}}", routes.GetVideoTranscription)
		})
	})

	return router
}
