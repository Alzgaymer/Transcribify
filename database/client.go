package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
	"transcribify/config"
)

type connectionFunc func() error

func NewClient(ctx context.Context, attemptsToConnect uint, sleep time.Duration) (client *pgx.Conn, err error) {
	configuration := config.DB()
	dsn := GetDSN(configuration)

	err = doWithAttempts(attemptsToConnect, sleep, func() error {

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

func doWithAttempts(attempts uint, sleep time.Duration, f connectionFunc) error {
	for i := 0; i < int(attempts); i++ {
		if err := f(); err != nil {
			return err
		}
		time.Sleep(sleep)
	}
	return nil
}
