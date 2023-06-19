package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"telegram_service/internal/config"
	"telegram_service/internal/server"
	"telegram_service/internal/service"
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

	authService := service.NewAuthService()

	tgService := service.TgService{}

	tgConnect := server.NewTelegram(&cfg, &tgService, authService)

	tgConnect.Start()
}
