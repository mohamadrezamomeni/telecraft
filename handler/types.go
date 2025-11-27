package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ResponseHandlerFunc struct {
	MessageConfigs []*tgbotapi.MessageConfig
	ReleaseState   bool
	RedirectRoot   bool
	Data           map[string]string
	Path           string
}

type Context struct {
	*tgbotapi.Update
	Data   map[string]any
	Params map[string]string
	UserID string
}

type HandlerFunc = func(*Context) (*ResponseHandlerFunc, error)

type Middleware = func(HandlerFunc) HandlerFunc
