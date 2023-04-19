package service

import (
	"log"
	"net/http"
	"os"
	"transcribify/pkg/auth"
	"transcribify/pkg/finders"
	"transcribify/pkg/repository"
)

type (
	Service struct {
		Manager       auth.TokenManager
		Authorization auth.Authorization
		client        *http.Client
		repo          repository.Video
	}
)

func New(repository repository.Repository, client *http.Client) *Service {
	manager, err := auth.NewManager(os.Getenv("JWT_SALT"))
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		Manager:       manager,
		Authorization: auth.NewAuthorizationManager(repository.User, manager),
		repo:          repository.Video,
		client:        client,
	}
}

// CacheVideoFinders
func (s *Service) CacheVideoFinders() []finders.Finder {
	return []finders.Finder{
		finders.NewDatabaseFinder(s.repo),
		finders.NewAPIFinder(s.client, s.repo),
	}
}
