package services

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type LatexClient struct{}

func (l *LatexClient) Test() {
	text := `
		\[
		\sum_{i=0}^\infty x_i = \sigma_x	
		\]
	`
	l.CreateImage(text)
}

func (l *LatexClient) Handler() tele.HandlerFunc {
	return func(ctx tele.Context) (err error) {
		go func() {
			input := ctx.Message().Payload
			l.CreateImage(input)
			b, err2 := os.ReadFile("latex.png")
			err = err2
			ctx.SendAlbum(tele.Album{&tele.Photo{
				File: tele.FromReader(bytes.NewReader(b)),
			}})
			fmt.Printf("\"salam\": %v\n", "salam")
		}()
		return
	}
}

func (l *LatexClient) CreateImage(text string) { // make this function work
	b, err := os.ReadFile("style/latex.html")
	if err != nil {
		log.Fatal(err)
	}
	html := fmt.Sprintf(string(b), l)
	file, err := os.Create("temp/latex/ltx.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte(html))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	err = exec.Command("./python html-to-image/main.py").Run()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func (l *LatexClient) Method() string {
	return "/latex"
}

func (l *LatexClient) Middleware() []tele.MiddlewareFunc {
	return []tele.MiddlewareFunc{
		middleware.Logger(),
		middleware.IgnoreVia(),
	}
}
