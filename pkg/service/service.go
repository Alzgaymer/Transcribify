package service

import (
	"log"
	"os"
	"transcribify/pkg/auth"
)

type (
	Service struct {
		Manager auth.TokenManager
	}
)

func New() *Service {
	manager, err := auth.NewManager(os.Getenv("JWT_SALT"))
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		Manager: manager,
	}
}
