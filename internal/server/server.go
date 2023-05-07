package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
	"transcribify/internal/routes"
	"transcribify/internal/routes/middlewares"
	"transcribify/pkg/dbclient"
	"transcribify/pkg/finders"
	"transcribify/pkg/hash"
	"transcribify/pkg/logging"
	repo "transcribify/pkg/repository"
	"transcribify/pkg/service"
)

func Server(ctx context.Context) *http.Server {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	repository := Repository(ctx)
	client := Client()
	return &http.Server{
		Addr: ":" + os.Getenv("APP_PORT"),
		Handler: Router(
			Logger(),
			client,
			service.New(*repository, finders.NewAPIFinder(client, repository.Video), hash.NewBCHasher(bcrypt.DefaultCost)),
			repository,
		),
	}
}

func Repository(ctx context.Context) *repo.Repository {

	client, err := dbclient.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return repo.NewRepositories(client, hash.NewBCHasher(bcrypt.DefaultCost))
}

// Router uses for http.Server struct Handler field.
// It implements several endpoints
func Router(logger *zap.Logger, client *http.Client, service *service.Services, repository *repo.Repository) http.Handler {

	router := chi.NewRouter()

	route := routes.NewRoute(
		logger, client, repository, service,
	)

	auth := middlewares.Identify(logger, service.Manager)

	router.Route("/api/v1", func(r chi.Router) {

		//GET /api/v1/video/{id}?lang=
		r.With(middlewares.LogVideoRequest(logger), auth).
			Get("/video/{id}", route.GetVideoTranscription)

		//GET /api/v1/user/history/{page}?limit=
		r.With(auth).
			Get("/user/history/{page}", route.GetUserVideo)

		//GET /api/v1/hello-world
		r.With(auth).
			Get("/hello-world", route.HelloWorld)

		r.Route("/auth", func(r chi.Router) {

			//POST /api/v1/auth/token
			r.Post("/token", route.GetToken)

			//POST /api/v1/auth/sign-up
			r.Post("/sign-up", route.SignUp)

			//POST /api/v1/auth/login
			r.Post("/login", route.LogIn)
		})
	})

	return router
}

func Logger() *zap.Logger {
	conf := zap.NewDevelopmentEncoderConfig()
	conf.EncodeTime = zapcore.TimeEncoderOfLayout(time.UnixDate)
	logg, err := logging.New(
		logging.WithDevelopment(true),
		logging.WithLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
		logging.WithEncoderConfig(conf),
	)
	if err != nil {
		log.Fatal(err)
	}

	return logg
}

// Client returns pointer to the http.Client with 30s timeout
func Client() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
