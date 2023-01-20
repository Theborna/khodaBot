package api

import (
	"kh-bot/api/services"

	tele "gopkg.in/telebot.v3"
)

type KhBot struct {
	*tele.Bot
}

// get an instance of the KhBot api, with additional services as a parameter
func NewBot(pref tele.Settings, services ...services.ServiceProvider) (*KhBot, error) {
	b, err := tele.NewBot(pref)
	Bot := &KhBot{
		Bot: b,
	}
	Bot.handlers(services...)
	return Bot, err
}

// add all function handles
func (bot *KhBot) handlers(services ...services.ServiceProvider) {
	for _, service := range services {
		bot.Handle(service.Method(), service.Handler(), service.Middleware()...)
	}
}
