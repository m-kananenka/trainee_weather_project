package main

import (
	"github.com/sirupsen/logrus"
	"trainee_weather_project/weather_service/internal/config"
	"trainee_weather_project/weather_service/internal/server"
)

func main() {
	cfg := config.Config{}
	logger := logrus.New()

	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	serv := server.NewWeatherServer(logger, &cfg)

	serv.PrintTemps()

}
