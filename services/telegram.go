package services

import (
	"github.com/biomaks/feederBot/settings"
	"github.com/yanzay/tbot"
)

type TelegramInterface interface {
	Send(text string, tm tbot.Message) (bool, error)
	Client() *tbot.Client
	Server() *tbot.Server
}

type TelegramService struct {
	client *tbot.Client
	server *tbot.Server
}

type Telegram struct {
	service TelegramInterface
}

func (t *TelegramService) Send(text string, tm tbot.Message) (bool, error) {
	_, err := t.client.SendMessage(text, tm.Chat.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TelegramService) Client() *tbot.Client {
	return t.client
}

func (t *TelegramService) Server() *tbot.Server {
	return t.server
}

func (t *Telegram) SendMessage(text string, tm tbot.Message) {
	_, err := t.service.Send(text, tm)
	if err != nil {
		panic(err)
	}
}

func (t *Telegram) GetClient() *tbot.Client {
	return t.service.Client()
}

func (t *Telegram) GetServer() *tbot.Server {
	return t.service.Server()
}

func NewTelegramService(token string) Telegram {
	config := settings.GetSettings()
	botSettings := config.Bot()
	bot := tbot.New(botSettings.Token)
	service := Telegram{&TelegramService{bot.Client(), bot}}
	return service
}
