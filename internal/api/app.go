package app

import (
	"Road_services/internal/app/dsn"
	"Road_services/internal/app/repository"

	"github.com/joho/godotenv"
)

type Application struct {
	repository *repository.Repository
}

func New() (Application, error) {
	_ = godotenv.Load()
	repo, err := repository.New(dsn.SetConnectionString())
	if err != nil {
		return Application{}, err
	}

	return Application{repository: repo}, nil
}
