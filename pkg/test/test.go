package test

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"time"
	"transcribify/internal/config"
	"transcribify/pkg/dbclient"
)

const PostgresVersion = "15"

func CreatePostgresContainer(seconds time.Duration, root string, migrateVersion uint) (error, func() error) {

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not construct pool: %w", err), nil
	}

	var timeout = seconds

	pool.MaxWait = timeout * time.Second

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancelFunc()

	err = pool.Client.PingWithContext(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to Docker: %w", err), nil
	}

	// pulls an image, creates a container based on it and runs it

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        PostgresVersion,
		Name:       "repository-test-postgres",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", os.Getenv("DB_PASSWORD")),
			fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("DB_USERNAME")),
			fmt.Sprintf("POSTGRES_DB=%s", os.Getenv("DB_DATABASE")),
			"listen_address='*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			docker.Port(os.Getenv("DB_PORT") + "/tcp"): {
				{
					HostIP:   os.Getenv("DB_HOST"),
					HostPort: os.Getenv("DB_PORT"),
				},
			},
		}
	})
	if err != nil {
		return fmt.Errorf("could not start container: %w", err), nil
	}
	databaseUrl := dbclient.GetDSN(config.DB())
	log.Println("Connecting to on url: ", databaseUrl)

	container.Expire(uint(timeout)) //Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet

	// Migrations using migrate package
	// Migrate if err == nil else panic

	if migrations, err := migrate.New(
		"file://"+root+"assets/migrations/postgres",
		databaseUrl); err == nil {
		if err := migrations.Migrate(migrateVersion); err != nil {
			return fmt.Errorf("couldn`t migrate to version %d:%w", migrateVersion, err), nil
		}
	} else {
		return err, nil
	}

	return nil, func() error {
		if err := pool.Purge(container); err != nil {
			return err
		}
		return nil
	}
}
