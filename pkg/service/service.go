package service

import (
	"log"
	"os"
	"transcribify/pkg/auth"
	"transcribify/pkg/repository"
)

type (
	Service struct {
		Manager       auth.TokenManager
		Authorization auth.Authorization
	}
)

func New(repository repository.Repository) *Service {
	manager, err := auth.NewManager(os.Getenv("JWT_SALT"))
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		Manager:       manager,
		Authorization: auth.NewAuthorizationManager(repository.User),
	}
}
