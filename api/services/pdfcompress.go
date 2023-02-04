package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type CompressorClient struct{}

const (
	HELP_COMPRESS             = `please reply to a file you want to compress`
	NO_FILE_COMPRESS          = `no media to compress`
	UNSUPPORTED_FILE_COMPRESS = `we only support documents`
	CAPTION                   = `compressed by @TheKhoda_bot`
)

func (p *CompressorClient) Handler() telebot.HandlerFunc {
	return func(ctx tele.Context) error {
		if !ctx.Message().IsReply() {
			return ctx.Send(HELP_COMPRESS)
		}
		m := ctx.Message().ReplyTo.Media()
		if m == nil {
			return ctx.Send(NO_FILE_COMPRESS)
		}
		name := "compressed"
		if len(ctx.Message().Payload) > 0 {
			name = ctx.Message().Payload
		}
		dpi := 160
		switch m.MediaType() {
		case "document":
			err := make(chan error)
			path := fmt.Sprintf("temp/pdf/%v.pdf", time.Now().Unix())
			start := time.Now()
			msg, _ := ctx.Bot().Send(ctx.Recipient(), "downloading file...")
			defer ctx.Bot().Delete(msg)
			go func() {
				e := ctx.Bot().Download(m.MediaFile(), path)
				err <- e
			}()
			select {
			case e := <-err:
				downloadTime := time.Since(start).Round(time.Millisecond)
				log.Printf("downloadTime: %v\n", downloadTime)
				ctx.Bot().Edit(msg, fmt.Sprintf("download successful.\ntook %v", downloadTime))
				if e != nil {
					return e
				}
				s, e := p.compressFile(path, dpi)
				if e != nil {
					return e
				}
				fmt.Printf("s: %v\n", s)
				defer os.Remove(s)
				file := &tele.Document{
					File:     tele.FromDisk(s),
					Caption:  CAPTION,
					FileName: fmt.Sprintf("%s.pdf", name),
				}
				_, e = file.Send(ctx.Bot(), ctx.Recipient(), &tele.SendOptions{ReplyTo: ctx.Message()})
				return e
			case <-time.After(5 * time.Second):
				return ctx.Reply("download timeout")
			}
		default:
			return ctx.Send(UNSUPPORTED_FILE_COMPRESS)
		}
	}
}
func (p *CompressorClient) compressFile(path string, dpi int) (string, error) {
	path2 := fmt.Sprintf("./temp/pdf/compressed_%v.pdf", time.Now().Unix())
	defer os.Remove(path)
	log.Printf("compressing file")
	ex := exec.Command("bash", "compress.sh", path, path2, fmt.Sprint(dpi))
	_, err := ex.Output()
	if err != nil {
		return "", err
	}
	log.Printf("compressed file")
	return path2, nil
}

func (p *CompressorClient) Method() string {
	return "/compress"
}
func (p CompressorClient) Middleware() []telebot.MiddlewareFunc {
	return []tele.MiddlewareFunc{
		// middleware.Logger(),
		middleware.IgnoreVia(),
	}
}
