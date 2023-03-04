package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"strings"
)

var (
	routeConfig = RouteConfiguration{
		filename: "route.env",
	}
	apiConfig = APIConfiguration{
		filename: "api.env",
	}
	dbConfig = DBConfiguration{
		filename: "db.env",
	}
)

func getProjectPath(inputPath string) (string, error) {
	index := strings.Index(inputPath, "yt-video-transcriptor")
	if index == -1 {
		return "", fmt.Errorf("project path not found in input path")
	}
	return inputPath[:index+len("yt-video-transcriptor")] + "\\", nil
}

func init() {
	wd, err := os.Getwd()
	checkErr(err)
	wd, err = getProjectPath(wd)
	checkErr(err)
	//route configuration
	err = cleanenv.ReadConfig(wd+routeConfig.filename, &routeConfig)
	checkErr(err)

	//api configuration
	err = cleanenv.ReadConfig(wd+apiConfig.filename, &apiConfig)
	checkErr(err)

	//db configuration
	err = cleanenv.ReadConfig(wd+dbConfig.filename, &dbConfig)
	checkErr(err)
}

func GetRoute() RouteConfiguration {
	return routeConfig
}

func GetAPI() APIConfiguration {
	return apiConfig
}

func GetDB() DBConfiguration {
	return dbConfig
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type RouteConfiguration struct {
	filename string
	Port     string `env:"APP_PORT"`
}

type APIConfiguration struct {
	filename string
	Key      string `env:"API_KEY"`
	API      string `env:"API_URL"`
}

type DBConfiguration struct {
	filename string
	Password string
	Host     string
	Port     string
}
