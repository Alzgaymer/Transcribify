package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	serve "transcribify/internal/server"
)

func Run() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())

	server := serve.Server(ctx)

	logger := serve.Logger()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Waiting to signal from os
	go func() {
		defer cancel()
		<-sig

		timeout, timeoutfn := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutfn()
		go func() {
			// Waiting until context is done then panic
			<-timeout.Done()
			if errors.Is(timeout.Err(), context.DeadlineExceeded) {
				logger.Fatal("Graceful shutdown timed out... forcing exit...")
			}
		}()

		err = server.Shutdown(timeout)
		if err != nil {
			logger.Fatal("Failed to shutdown server", zap.Error(err))
		}
	}()

	logger.Info("Server is running", zap.String("port", os.Getenv("APP_PORT")))

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to serve", zap.Error(err))
	}

	<-ctx.Done()
}
