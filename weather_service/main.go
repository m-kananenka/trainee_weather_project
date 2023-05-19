package main

import (
	"github.com/sirupsen/logrus"
	"weather_service/internal/config"
	"weather_service/internal/server"
	"weather_service/internal/service"
)

func main() {
	cfg := config.Config{}
	logger := logrus.New()

	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	service := service.NewGRPCServer(&cfg, logger)

	serv := server.NewWeatherServer(logger, &cfg, service)

	serv.Register()
	serv.Start()

}
