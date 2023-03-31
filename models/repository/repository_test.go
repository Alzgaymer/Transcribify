package repository

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
	"yt-video-transcriptor/config"
	"yt-video-transcriptor/database"
	"yt-video-transcriptor/models"
)

func Test_formatQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "With \\n",
			query:    "\nCREATE DATABASE example;",
			expected: " CREATE DATABASE example;",
		},
		{
			name:     "With \\t",
			query:    "\tCREATE DATABASE example;",
			expected: "CREATE DATABASE example;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatQuery(tt.query)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

var db *pgx.Conn

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", os.Getenv("DB_PASSWORD")),
			fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("DB_USERNAME")),
			fmt.Sprintf("POSTGRES_DB=%s", os.Getenv("DB_DATABASE")),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	databaseUrl := database.GetDSN(config.DB())

	log.Println("Connecting to database on url: ", databaseUrl)

	var timeout time.Duration = 120
	container.Expire(uint(timeout)) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = timeout * time.Second

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	if err = pool.Retry(func() error {
		db, err = pgx.Connect(ctx, databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Migrations using migrate package
	migrations, err := migrate.New(
		"file://migrations/postgres",
		databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := migrations.Migrate(1); err != nil {
		log.Fatal(err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}

func TestYTVideoRepositoryRepository(t *testing.T) {
	repo := NewYTVideoRepository(db)
	type testFunc func(
		ctx context.Context,
		repository *YTVideoRepository,
		video models.YTVideo,
		request models.VideoRequest,
		id int,
	) (int, models.YTVideo, error)

	type testStruct struct {
		name string

		id         int
		returnedId int

		model         models.YTVideo
		returnedModel models.YTVideo

		request       models.VideoRequest
		Do            testFunc
		expectedError error
	}
	testData := []testStruct{
		{
			name: "Successful Create",
			model: models.YTVideo{
				Transcription: []models.Transcription{
					{
						Subtitle: "test",
					},
				},
			},
			returnedModel: models.YTVideo{},

			id:         1,
			returnedId: 1,

			request: models.VideoRequest{
				VideoID:  "00000000000",
				Language: "ua",
			},
			Do: func(ctx context.Context, repository *YTVideoRepository, video models.YTVideo, request models.VideoRequest, id int) (int, models.YTVideo, error) {
				id, err := repository.Create(ctx, video, request)
				return id, models.YTVideo{}, err
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {

		})
	}
}
