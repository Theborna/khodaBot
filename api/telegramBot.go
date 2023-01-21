package api

import (
	"fmt"
	"kh-bot/api/services"
	"strings"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

type KhBot struct {
	*tele.Bot
}

var (
	WELCOME_TEXT = `
	Hi!
	I'm Khoda, I'm a modular telegram bot with helps you work with
	various other services through telegram!
	I'm maintained by https://t.me/bornaKhodabandeh and you can see my source code
	at https://github.com/Theborna/khodaBot
	`
	CURRENT_APIS = []string{"wolfram alpha", "latex to image", "gpt3 client", "code to image", "stable diffusion"}
)

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
	bot.Handle("/start", startHandler())
	for _, service := range services {
		bot.Handle(service.Method(), service.Handler(), service.Middleware()...)
	}
}

func startHandler() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		return ctx.Send(WELCOME_TEXT,
			fmt.Sprintf("current api's:\n%v", strings.Join(CURRENT_APIS, ", ")))
	}
}
