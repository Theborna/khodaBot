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
	// wolframClient, err := services.NewWolframClient(os.Getenv("WOLFRAM_APP_ID"))
	// errHandler(err)
	gptClient, err := services.NewGptClient(os.Getenv("CHAT_GPT_APP_ID"))
	errHandler(err)
	latexClient := &services.LatexClient{}
	go latexClient.Test()
	// codeClient := &services.CodeClient{}
	// go codeClient.Test()
	stableDiffClient := &services.SDClient{}
	// _, err := stableDiffClient.GetLink(1, []string{"woman, beautiful, elegany, golden braided hair golden eyes, green and gold caftan, dress, smiling, fairy, shiny, realistic ,4k"}, []string{})
	errHandler(err)
	Bot, err := api.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_API"),
		Poller: Poller,
	}, gptClient, latexClient, stableDiffClient)
	// })
	errHandler(err)
	Bot.Start()
}

func errHandler(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
