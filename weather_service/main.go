package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"weather_service/internal/config"
	"weather_service/internal/server"
	"weather_service/internal/service"
)

func main() {
	cfg := config.Config{}
	logger := logrus.New()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	err = cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	service := service.NewGRPCServer(&cfg, logger)

	serv := server.NewWeatherServer(logger, &cfg, service)

	serv.Register()
	serv.Start()

}
