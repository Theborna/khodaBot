package api

import (
	"fmt"
	"kh-bot/api/services"
	"log"
	"os"
	"strings"
	"time"

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
	CURRENT_APIS = []string{"wolfram alpha", "latex to image", "gpt3 client", "code to image", "stable diffusion", "pdf compressor client"}
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
	bot.Handle("/hi", startHandler())
	bot.Handle("/help", startHandler())
	bot.Handle("/ping", bot.pongHandler())
	bot.Handle("/report", bot.reportHandler())
	for _, service := range services {
		bot.Handle(service.Method(), service.Handler(), service.Middleware()...)
	}
}

func startHandler() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		ctx.Send(WELCOME_TEXT)
		return ctx.Send(fmt.Sprintf("current api's:\n%v", strings.Join(CURRENT_APIS, ", ")))
	}
}
func (bot *KhBot) pongHandler() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		start := time.Now()
		myMsg := make(chan telebot.Message)
		go func() {
			msg, _ := bot.Reply(ctx.Message(), "Pinging...")
			myMsg <- *msg
		}()
		resp := <-myMsg
		ping := time.Since(start).Round(time.Millisecond)
		_, err := bot.Edit(&resp, fmt.Sprintf("PONG!!!\n%v", ping))
		return err
	}
}

func (b *KhBot) reportHandler() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		sender, message := ctx.Sender().Username, ctx.Message().Payload
		text := fmt.Sprintf("sender: t.me/%s, message: %s\n", sender, message)
		if len(message) > 0 {
			path := fmt.Sprintf("./reports/report_%v.txt", time.Now().Format("01-02-2006"))
			f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Print(err)
				return ctx.Send("failed to write report")
			}
			defer f.Close()
			if _, err = f.WriteString(text); err != nil {
				log.Print(err)
				return ctx.Send("failed to write report")
			}
		}
		return ctx.Send("report sent successfully")
	}
}
