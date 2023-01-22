package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type CodeClient struct{}

func (l *CodeClient) Test() {
	text := `
	package main
	
	import (
		"fmt"
		"net/http"
		"time"
	)
	
	func greet(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! %s", time.Now())
	}
	
	func main() {
		http.HandleFunc("/", greet)
		http.ListenAndServe(":8080", nil)
	}
	`
	l.CreateImage(text)
}

func (l *CodeClient) Handler() tele.HandlerFunc {
	return func(ctx tele.Context) (err error) {
		msg, _ := ctx.Bot().Reply(ctx.Message(), "getting results...")
		go func() {
			defer ctx.Bot().Delete(msg)
			input := ctx.Message().Text[len(l.Method()):]
			if ctx.Message().IsReply() {
				input = ctx.Message().ReplyTo.Text
			}
			fmt.Printf("input: %v\n", input)
			l.CreateImage(input)
			b, err2 := os.ReadFile("./temp/code/tmp.png")
			err = err2
			ctx.SendAlbum(tele.Album{&tele.Photo{
				File: tele.FromReader(bytes.NewReader(b)),
			}})
			fmt.Printf("\"salam\": %v\n", "salam")
		}()
		return
	}
}

func (l *CodeClient) CreateImage(text string) {
	b, err := os.ReadFile("style/code.html")
	errHandler(err)
	html := fmt.Sprintf(string(b), text)
	file, err := os.Create("temp/code/code.html")
	errHandler(err)
	defer file.Close()
	file.Write([]byte(html))
	errHandler(err)
	err = exec.Command("bash", "screenshot.sh", "./temp/code/code.html",
		"./temp/code/tmp.png").Run()
	errHandler(err)
}

func (l *CodeClient) Method() string {
	return "/code"
}

func (l *CodeClient) Middleware() []tele.MiddlewareFunc {
	return []tele.MiddlewareFunc{
		middleware.Logger(),
		middleware.IgnoreVia(),
	}
}
