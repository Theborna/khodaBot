package services

import "gopkg.in/telebot.v3"

/*
any type of command the bot can execute as a microservice
*/ /*
should be used as a wrapper that wraps any api as a telegram command
*/
type ServiceProvider interface {
	Handler() telebot.HandlerFunc
	Method() string
	Middleware() []telebot.MiddlewareFunc
}
