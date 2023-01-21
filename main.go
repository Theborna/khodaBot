package main

import (
	"fmt"
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
	// gptClient, err := services.NewGptClient(os.Getenv("CHAT_GPT_APP_ID"))
	// errHandler(err)
	// latexClient := &services.LatexClient{}
	// go latexClient.Test()
	// codeClient := &services.CodeClient{}
	// go codeClient.Test()
	stableDiffClient, err := services.NewSDClient(os.Getenv("REPLICATE_API_KEY"))
	errHandler(err)
	fmt.Printf("stableDiffClient.Request: %v\n", stableDiffClient.Request)
	err = stableDiffClient.Create()
	errHandler(err)
	fmt.Printf("stableDiffClient.Response: %v\n", stableDiffClient.Response)
	// Bot, err := api.NewBot(tele.Settings{
	// 	Token:  os.Getenv("TELEGRAM_API"),
	// 	Poller: Poller,
	// }, latexClient, codeClient)
	// errHandler(err)
	// Bot.Start()
}

func errHandler(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
