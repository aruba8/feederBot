package services

import (
	"github.com/biomaks/feederBot/settings"
	"github.com/yanzay/tbot/v2"
)

type TelegramInterface interface {
	Send(text string, chatID string) (bool, error)
	Client() *tbot.Client
	Server() *tbot.Server
}

type TelegramService struct {
	client *tbot.Client
	server *tbot.Server
}

type Telegram struct {
	Service TelegramInterface
}

func (t *TelegramService) Send(chatID string, text string) (bool, error) {
	_, err := t.client.SendMessage(chatID, text)
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

func (t *Telegram) SendMessage(chatID string, text string) {
	_, err := t.Service.Send(chatID, text)
	if err != nil {
		panic(err)
	}
}

func (t *Telegram) GetClient() *tbot.Client {
	return t.Service.Client()
}

func (t *Telegram) GetServer() *tbot.Server {
	return t.Service.Server()
}

func NewTelegramService() Telegram {
	config := settings.GetSettings()
	botSettings := config.Bot()
	bot := tbot.New(botSettings.Token)
	service := Telegram{&TelegramService{bot.Client(), bot}}
	return service
}
