package app

import (
	"time"

	"Road_services/internal/app/dsn"
	"Road_services/internal/app/redis"
	"Road_services/internal/app/repository"

	// "Road_services/internal/app/controllers"
	"github.com/golang-jwt/jwt"

	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

type Application struct {
	config     *Config
	repository *repository.Repository
	redis      *redis.Client
}

type Config struct {
	JWT struct {
		Token         string
		SigningMethod jwt.SigningMethod
		ExpiresIn     time.Duration
	}
}

func New() (*Application, error) {
	_ = godotenv.Load()

	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(dsn.SetConnectionString())
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New()
	if err != nil {
		return nil, err
	}

	return &Application{config: config, repository: repo, redis: redisClient}, nil
}
