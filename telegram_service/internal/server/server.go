package server

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegram_service/internal/config"
	"telegram_service/internal/service"
)

type Telegram struct {
	cfg         *config.Config
	tgService   *service.TgService
	authService *service.AuthService
}

func NewTelegram(cfg *config.Config, tgService *service.TgService, auth *service.AuthService) Telegram {
	return Telegram{
		cfg:         cfg,
		tgService:   tgService,
		authService: auth,
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

			var message string
			if t.authService.CheckAuth(update.Message.Chat.ID) {
				message = t.GetWeather(update, err)
			} else {
				if t.authService.Auth(update.Message.Text, update.Message.Chat.ID) {
					message = "You are successfully authorized. \nSelect the city where you want to know the weather. " +
						"\nFor example: Minsk"
				} else {
					message = "Write correct login and password, please"
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func (t *Telegram) GetWeather(update tgbotapi.Update, err error) string {
	var message string

	if update.Message.Text != "" {
		message, err = t.tgService.GetWeather(update.Message.Text)
		if err != nil || message == "" {
			message = "Incorrect input"
		}

	} else {
		message = "Write city please"
	}
	return message
}
