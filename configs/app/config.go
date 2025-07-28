package appConfig

import (
	"github.com/joho/godotenv"
	"os"
)

type App struct {
	Name string
	Env  string
}

func New() (*App, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}, nil
}
