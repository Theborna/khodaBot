package services

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"

	wolfram "github.com/Krognol/go-wolfram"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const PROMPT_LEN_LIMIT = 5

type WolframClient struct {
	*wolfram.Client
}

func NewWolframClient(appID string) (*WolframClient, error) {
	client := wolfram.Client{AppID: appID}
	simpleQuery := "what is wolfram alpha"
	ans, err := client.GetShortAnswerQuery(simpleQuery, wolfram.Metric, 100)
	fmt.Printf("%v: %v\n", simpleQuery, ans)
	return &WolframClient{Client: &client}, err
}

func (m *WolframClient) Method() string {
	return "/wolfram"
}

func (m *WolframClient) Handler() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		query := ctx.Message().Payload
		if len(query) < PROMPT_LEN_LIMIT {
			return ctx.Reply("too short!")
		}
		wg := sync.WaitGroup{}
		var res string
		var fail error
		fail = nil
		errHandler := func(err error) {
			if err != nil {
				log.Println(err.Error())
				ctx.Reply(fmt.Sprintf("unexpected error: %v", err.Error()))
				fail = err
			}
		}
		msg, _ := ctx.Bot().Reply(ctx.Message(), "getting results...")
		defer ctx.Bot().Delete(msg)
		wg.Add(1)
		go func() {
			log.Println("getting spoken answer")
			resp, err := m.GetSpokentAnswerQuery(query, wolfram.Metric, 2)
			errHandler(err)
			wg.Done()
			log.Println("got spoken answer")
			res = resp
		}()
		var res2 io.ReadCloser
		u := url.Values{}
		u.Add("timeout", "2")
		u.Add("units", "metric")
		wg.Add(1)
		go func() {
			log.Println("getting image answer")
			resp, q, err := m.GetSimpleQuery(query, u)
			errHandler(err)
			log.Printf("q: %v\n", q)
			res2 = resp
			wg.Done()
			log.Println("got image answer")
		}()
		wg.Wait()
		if fail != nil {
			return fail
		}
		log.Println("sent result")
		return ctx.SendAlbum(tele.Album{&tele.Photo{
			File:    tele.FromReader(res2),
			Caption: res,
		}})
	}
}

func (m *WolframClient) Middleware() []tele.MiddlewareFunc {
	return []tele.MiddlewareFunc{
		middleware.Logger(),
		middleware.IgnoreVia(),
	}
}
