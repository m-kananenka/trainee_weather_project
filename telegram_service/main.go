package main

import (
	"github.com/sirupsen/logrus"
	"telegram_service/internal/config"
	"telegram_service/internal/server"
	"telegram_service/internal/service"
)

func main() {
	cfg := config.Config{}

	logger := logrus.New()

	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	tgService := service.TgService{}

	tgConnect := server.NewTelegram(&cfg, &tgService)

	tgConnect.Start()
}
