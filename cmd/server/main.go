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
	var configuration config.AppConfiguration
	err := cleanenv.ReadConfig("route.env", configuration)
	if err != nil {
		panic(err)
	}

	server := http.Server{Addr: ":" + configuration.Port, Handler: service()}
	log, err := logger.New(
		logger.WithDevelopment(true),
		logger.WithEncoding("console"),
		logger.WithLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
	)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer cancel()
		<-sig

		timeout, _ := context.WithTimeout(ctx, 30*time.Second)

		go func() {
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

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to serve", zap.Error(err))
	}
	<-ctx.Done()
}

func service() http.Handler {

	router := chi.NewRouter()

	router.Get("/api/v1/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello"))
	})

	return router
}
