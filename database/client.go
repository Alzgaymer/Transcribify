package database

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"time"
	"transcribify/config"
)

func NewClient(ctx context.Context) (client *pgx.Conn, err error) {
	configuration := config.DB()
	dsn := GetDSN(configuration)

	err = doWithAttempts(ctx, func() error {

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		client, err = pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
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
	return backoff.Retry(operation, backOff)
}
