package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const (
	version = "f178fa7a1ae43a9a9af01b833b9d2ecf97b1bcb0acfd2dc5dd04895e042863f1"
	owner   = "stability-ai"
	name    = "stable-diffusion"
)

type SDClient struct{}
type resultUrl string

const HELP = `
the stable-diffusion method works by using the stable-diffusion api from
https://replicate.com/stability-ai/stable-diffusion and sending it via telegram

write a valid stable diffusion query and then reply to it with /stable_diff

example: 

p: An astronaut riding a horse in photorealistic style
n: planet earth
k: 3

this will generate k pictures of 'An astronaut riding a horse in photorealistic style'
without 'planet earth' inside it
`
const PY_ERR = `exit status 1`

func (c SDClient) Handler() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		if !ctx.Message().IsReply() {
			return ctx.Send(HELP)
		}
		input := ctx.Message().ReplyTo.Text
		n, prompt, negPrompt, err := c.parseInput(input)
		if err != nil {
			return ctx.Send(err.Error())
		}
		resultUrl, err := c.GetLink(n, prompt, negPrompt)
		if err != nil {
			if err.Error() == PY_ERR {
				err = fmt.Errorf("api failed to respond")
			}
			return ctx.Send(err.Error())
		}
		var album tele.Album
		for _, u := range resultUrl {
			album = append(album, &tele.Photo{
				File: tele.FromURL(string(u)),
			})
		}
		return ctx.SendAlbum(album)
	}
}

func (c SDClient) GetLink(n int, prompt, negPrompt string) ([]resultUrl, error) {
	input := struct {
		Name       string `json:"name"`
		Owner      string `json:"owner"`
		Version    string `json:"version"`
		Prompt     string `json:"prompt"`
		NegPrompt  string `json:"neg_prompt"`
		NumOutputs int    `json:"num_outputs"`
	}{
		Name:       name,
		Owner:      owner,
		Version:    version,
		Prompt:     prompt,
		NegPrompt:  negPrompt,
		NumOutputs: 1,
	}
	inputJson, _ := json.Marshal(input)
	path := fmt.Sprintf("temp/diff/tmp%d.json", time.Now().Unix())
	file, err := os.Create(path)
	errHandler(err)
	defer file.Close()
	file.Write(inputJson)
	var res []resultUrl
	output := make(chan resultUrl)
	go func() {
		wg := &sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cmd := exec.Command("bash", "stable_diff.sh", path)
				out, e := cmd.Output()
				fmt.Printf("out: %v\n", out)
				if err == nil && e != nil {
					err = e
					return
				}
				var result []resultUrl
				e2 := json.Unmarshal(out, &result)
				if err == nil && e2 != nil {
					err = e2
					return
				}
				fmt.Printf("result: %v\n", result)
				if len(result) > 0 {
					output <- result[0]
				}
			}()
		}
		wg.Wait()
		close(output)
	}()
	for s := range output {
		res = append(res, s)
	}
	fmt.Printf("len(res): %v\n", len(res))
	return res, err
}

func (c SDClient) parseInput(input string) (int, string, string, error) {
	input = strings.ReplaceAll(input, "\n", " ") // remove all new lines
	input = strings.ReplaceAll(input, "###", " ")
	keywords := []string{"k:", "p:", "n:"}
	for _, key := range keywords {
		m, _ := regexp.Compile(key)
		input = m.ReplaceAllString(input, fmt.Sprintf("###%s", key))
	}
	inputs := strings.Split(input, "###")
	result := map[string]string{"positive": "", "negative": ""}
	n := 1 // default value
	for _, i := range inputs {
		if len(i) < 1 {
			continue
		}
		switch i[0] {
		case 'p':
			result["positive"] += strings.Trim(i[strings.Index(i, ":")+1:], " ") + ", "
		case 'n':
			result["negative"] += strings.Trim(i[strings.Index(i, ":")+1:], " ") + ", "
		case 'k':
			m, _ := regexp.Compile(`\d+`)
			n, _ = strconv.Atoi(string(m.Find([]byte(i))))
		}
	}
	if n > 4 {
		n = 4
	}
	return n, result["positive"], result["negative"], nil
}

func (c SDClient) Method() string {
	return "/stable_diff"
}

func (c *SDClient) Middleware() []tele.MiddlewareFunc {
	return []tele.MiddlewareFunc{
		middleware.Logger(),
		middleware.IgnoreVia(),
	}
}
