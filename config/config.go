package config

import (
	"os"
)

func GetRoute() RouteConfiguration {
	return RouteConfiguration{
		Port: os.Getenv("APP_PORT"),
	}
}
func GetAPI() APIConfiguration {
	return APIConfiguration{
		Key: os.Getenv("API_KEY"),
		API: os.Getenv("API_URL"),
	}
}

func GetDB() DBConfiguration {
	return DBConfiguration{}
}

type RouteConfiguration struct {
	Port string `env:"APP_PORT"`
}

type APIConfiguration struct {
	Key string `env:"API_KEY"`
	API string `env:"API_URL"`
}

type DBConfiguration struct {
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
}
