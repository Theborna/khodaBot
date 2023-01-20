package main

import (
	"kh-bot/api"
	"kh-bot/api/services"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var (
	Poller = &tele.LongPoller{Timeout: 10 * time.Second}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		panic("Error loading .env file")
	}
	log.Printf("loaded environment variables")
}

func main() {
	wolframClient, err := services.NewWolframClient(os.Getenv("WOLFRAM_APP_ID"))
	errHandler(err)
	gptClient, err := services.NewGptClient(os.Getenv("CHAT_GPT_APP_ID"))
	errHandler(err)
	Bot, err := api.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_API"),
		Poller: Poller,
	}, wolframClient, &services.LatexClient{}, gptClient)
	errHandler(err)
	Bot.Start()
}

func errHandler(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
