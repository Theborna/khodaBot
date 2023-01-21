package services

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type LatexClient struct{}

var errHandler = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
			input := ctx.Message().ReplyTo.Text[len(l.Method()):]
			if ctx.Message().IsReply() {
				input = ctx.Message().ReplyTo.Text
			}
			input = fixInput(input)
			fmt.Printf("input: %v\n", input)
			l.CreateImage(input)
			b, err2 := os.ReadFile("./temp/latex/tmp.png")
			err = err2
			ctx.SendAlbum(tele.Album{&tele.Photo{
				File: tele.FromReader(bytes.NewReader(b)),
			}})
			fmt.Printf("\"salam\": %v\n", "salam")
		}()
		return
	}
}

func fixInput(s string) string {
	m, _ := regexp.Compile("/")
	return m.ReplaceAllString(s, "\\")
}

func (l *LatexClient) CreateImage(text string) {
	b, err := os.ReadFile("style/latex.html")
	errHandler(err)
	html := fmt.Sprintf(string(b), text)
	file, err := os.Create("temp/latex/ltx.html")
	errHandler(err)
	defer file.Close()
	file.Write([]byte(html))
	errHandler(err)
	numLines := strings.Count(text, "\n") - 3
	err = exec.Command("bash", "screenshot.sh", "./temp/latex/ltx.html",
		"./temp/latex", "400", strconv.Itoa(300+numLines*10)).Run()
	errHandler(err)
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
