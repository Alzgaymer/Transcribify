package config

import (
	_ "github.com/joho/godotenv/autoload"
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
	Password string
	Host     string
	Port     string
}
