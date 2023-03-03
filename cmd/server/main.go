package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	log, err := logger.New()
	if err != nil {
		return nil
	}

	router := chi.NewRouter()

	// Create a route for the GET method that accepts the video ID as a parameter
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/{videoID:[a-zA-Z0-9_-]{11}}", func(r chi.Router) {
			r.Get("/{language:[a-zA-Z]{2}}", func(w http.ResponseWriter, r *http.Request) {
				// Handle GET request for video with specified language
				// Get the video ID from the URL parameters
				videoID := chi.URLParam(r, "videoID")
				language := chi.URLParam(r, "language")

				// Do something with the video ID, for example print it to the console
				log.Info("Input parameter",
					zap.String("Video ID", videoID),
					zap.String("Language", language),
					zap.String("URL", r.URL.Path),
				)

				var configuration config.APIConfiguration
				err := cleanenv.ReadConfig("api.env", &configuration)
				if err != nil {
					log.Error("Failed to read api.env", zap.Error(err))
				}

				// Return a response to the client
				_, err = w.Write([]byte("Video ID: " + videoID + "\n"))
				if err != nil {
					log.Error("Failed to write", zap.Error(err))
				}

				url := fmt.Sprintf("https://youtube-transcriptor.p.rapidapi.com/transcript?video_id=%s&lang=%s", videoID, language)

				req, _ := http.NewRequest(http.MethodGet, url, nil)

				req.Header.Add("X-RapidAPI-Key", configuration.Key)
				req.Header.Add("X-RapidAPI-Host", configuration.API)

				res, _ := http.DefaultClient.Do(req)

				defer res.Body.Close()

				log.Info("Getting response...")

				// Use the decoder to parse the response
				var data []config.YTVideo
				err = json.NewDecoder(res.Body).Decode(&data)
				if err != nil {
					log.Error("Failed to unmarshal data", zap.Error(err))
					w.WriteHeader(http.StatusConflict)
					return
				}

				w.WriteHeader(http.StatusOK)
				encoder := json.NewEncoder(w)
				encoder.SetIndent("", "    ")
				encoder.Encode(data)
			})
		})
	})

	return router
}
