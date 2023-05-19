package server

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegram_service/internal/config"
	"telegram_service/internal/service"
)

type Telegram struct {
	cfg       *config.Config
	tgService *service.TgService
}

func NewTelegram(cfg *config.Config, tgService *service.TgService) Telegram {
	return Telegram{
		cfg:       cfg,
		tgService: tgService,
	}
}

func (t *Telegram) Start() {
	bot, err := tgbotapi.NewBotAPI(t.cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			weather := t.getWeather()

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func (t *Telegram) getWeather() string {
	weather := t.tgService.GetWeather()

	return weather
}
