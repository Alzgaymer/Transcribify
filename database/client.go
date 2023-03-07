package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
	"yt-video-transcriptor/config"
)

func NewClient(ctx context.Context, attemptsToConnect uint, sleep time.Duration) (client Client, err error) {
	configuration := config.GetDB()
	dsn := getDSN(configuration)

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

func getDSN(configuration config.DBConfiguration) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		configuration.Username,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.Database,
	)
}

type connectionFunc func() error

func doWithAttempts(attempts uint, sleep time.Duration, f connectionFunc) error {
	for i := 0; i < int(attempts); i++ {
		if err := f(); err != nil {
			return err
		}
		time.Sleep(sleep)
	}
	return nil
}
