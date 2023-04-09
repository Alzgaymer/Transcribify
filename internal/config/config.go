package config

import (
	"os"
)

func Route() RouteConfiguration {
	return RouteConfiguration{
		Port: os.Getenv("APP_PORT"),
	}
}

func DB() DBConfiguration {
	return DBConfiguration{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}
}

type RouteConfiguration struct {
	Port string `env:"APP_PORT"`
}

type DBConfiguration struct {
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
}
