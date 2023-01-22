package services

import (
	"context"
	"fmt"
	"log"
	"time"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type GptCLient struct {
	gpt3.Client
}

func NewGptClient(apiKey string) (*GptCLient, error) {
	c := gpt3.NewClient(apiKey)
	client := &GptCLient{
		Client: c,
	}
	ctx := context.Background()
	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{"GPT 3 is"},
		MaxTokens: gpt3.IntPtr(30),
		Stop:      []string{"."},
		Echo:      true,
	})
	if err != nil {
		return client, nil
	}
	fmt.Println(resp.Choices[0].Text)
	return client, nil
}

func (g *GptCLient) Handler() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		res := make(chan string)
		log.Println("getting gpt result...")
		msg, _ := ctx.Bot().Reply(ctx.Message(), "getting results...")
		go func() {
			defer ctx.Bot().Delete(msg)
			defer close(res)
			msg := ctx.Message().Payload
			resp, err := g.Completion(context.Background(), gpt3.CompletionRequest{
				Prompt:    []string{msg},
				MaxTokens: gpt3.IntPtr(30),
				Stop:      []string{"."},
				Echo:      true,
			})
			if err != nil {
				res <- err.Error()
			} else {
				res <- resp.Choices[0].Text
			}
			log.Println("got gpt result")
		}()
		select {
		case result := <-res:
			fmt.Printf("sent gpt result: %s\n", result)
			return ctx.Reply(result)
		case <-time.After(1000 * time.Millisecond):
			result := "timed out waiting for gpt api"
			fmt.Println(result)
			return ctx.Reply(result)
		}
	}
}

func (g *GptCLient) Method() string {
	return "/gpt"
}

func (g *GptCLient) Middleware() []tele.MiddlewareFunc {
	return []tele.MiddlewareFunc{
		middleware.Logger(),
		middleware.IgnoreVia(),
	}
}
