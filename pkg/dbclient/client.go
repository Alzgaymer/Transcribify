package dbclient

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"
	"transcribify/internal/config"
	"transcribify/pkg/logging"
)

func NewClient(ctx context.Context) (client *pgx.Conn, err error) {
	configuration := config.DB()
	dsn := GetDSN(configuration)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = doWithAttempts(ctx, func() error {

		client, err = pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return client.Ping(ctx)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetDSN Uses fmt.Sprintf no need to test
func GetDSN(configuration config.DBConfiguration) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		configuration.Username,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.Database,
	)
}

func doWithAttempts(ctx context.Context, operation backoff.Operation) error {
	backOff := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	logger, err := logging.New(
		logging.WithOutputPaths("stderr"),
	)
	if err != nil {
		return err
	}
	return backoff.RetryNotify(operation, backOff,
		func(err error, duration time.Duration) {
			logger.Info("Connecting to dbclient",
				zap.Error(err),
				zap.Duration("Waiting for", duration),
			)
		})
}
