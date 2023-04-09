package repository

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
	"transcribify/internal/config"
	"transcribify/internal/models"
	"transcribify/pkg/dbclient"
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

const (
	PathToRoot      = "../../"
	MigrateVersion  = 1
	PostgresVersion = "15"
)

func TestMain(m *testing.M) {

	err := godotenv.Load(PathToRoot + ".env")
	if err != nil {
		log.Fatal(err)
	}
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	var timeout time.Duration = 120

	pool.MaxWait = timeout * time.Second

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancelFunc()

	err = pool.Client.PingWithContext(ctx)
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
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
		log.Fatalf("Could not start container: %s", err)
	}

	databaseUrl := dbclient.GetDSN(config.DB())

	log.Println("Connecting to dbclient on url: ", databaseUrl)

	container.Expire(uint(timeout)) //Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet

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
	// Migrate if err == nil else panic
	if migrations, err := migrate.New(
		"file://"+PathToRoot+"internal/migrations/postgres",
		databaseUrl); err == nil {
		if err := migrations.Migrate(MigrateVersion); err != nil {
			log.Fatal(err)
		}
	} else {
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

func TestYTVideoRepository(t *testing.T) {
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
		expectedId int

		video         models.YTVideo
		expectedVideo models.YTVideo

		request       models.VideoRequest
		Do            testFunc
		expectedError error
	}
	testData := []testStruct{
		{
			name: "Successful Create",
			video: models.YTVideo{
				Transcription: []models.Transcription{
					{
						Subtitle: "test",
					},
				},
			},
			expectedVideo: models.YTVideo{},

			id:         1,
			expectedId: 1,

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
		{
			name:  "Successful Read",
			video: models.YTVideo{},
			expectedVideo: models.YTVideo{
				Transcription: []models.Transcription{
					{
						Subtitle: "test",
					},
				},
			},

			id:         1,
			expectedId: 1,

			request: models.VideoRequest{
				VideoID:  "00000000000",
				Language: "ua",
			},
			Do: func(ctx context.Context, repository *YTVideoRepository, video models.YTVideo, request models.VideoRequest, id int) (int, models.YTVideo, error) {
				videoFromRead, err := repository.Read(ctx, request)
				return id, videoFromRead, err

			},
			expectedError: nil,
		},
		{
			name: "Successful Update",
			video: models.YTVideo{
				Transcription: []models.Transcription{
					{
						Subtitle: "test1",
					},
				},
			},
			expectedVideo: models.YTVideo{
				Transcription: []models.Transcription{
					{
						Subtitle: "test1",
					},
				},
			},

			id:         1,
			expectedId: 1,

			request: models.VideoRequest{
				VideoID:  "00000000000",
				Language: "ua",
			},
			Do: func(ctx context.Context, repository *YTVideoRepository, video models.YTVideo, request models.VideoRequest, id int) (int, models.YTVideo, error) {
				err := repository.Update(ctx, request, video)
				return id, video, err

			},
			expectedError: nil,
		},
		{
			name: "Successful Delete", video: models.YTVideo{}, expectedVideo: models.YTVideo{},
			id: 1, expectedId: 1,

			request: models.VideoRequest{
				VideoID:  "00000000000",
				Language: "ua",
			},
			Do: func(ctx context.Context, repository *YTVideoRepository, video models.YTVideo, request models.VideoRequest, id int) (int, models.YTVideo, error) {
				err := repository.Delete(ctx, request)
				return id, video, err

			},
			expectedError: nil,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			id, video, err := testCase.Do(ctx, repo,
				testCase.video,
				testCase.request,
				testCase.id,
			)

			assert.Equal(t, testCase.expectedId, id)
			assert.Equal(t, testCase.expectedVideo, video)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
