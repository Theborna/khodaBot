# khoda bot

Hi!
I'm Khoda, I'm a modular telegram bot written in `golang` with helps you work with
various other services through telegram!
I'm maintained by [borna in telegram](https://t.me/bornaKhodabandeh)

## current api's

----------------------------------------------------------------

- wolfram alpha
- latex to image
- gpt3 client
- code to image
- stable diffusion

## modules

----------------------------------------------------------------

services are added to the bot using the `serviceProvider` interface.

```golang
type ServiceProvider interface {
	Handler() telebot.HandlerFunc
	Method() string
	Middleware() []telebot.MiddlewareFunc
}
```

Each `serviceProvider` is simply a wrapper around an api
which converts the api in a way that can be used as a telegram
method.