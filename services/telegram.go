package services

import (
	"github.com/biomaks/feederBot/settings"
	"github.com/yanzay/tbot/v2"
	"log"
)

type BotClientInterface interface {
	SendMessage(chatID string, text string) error
}

type telegramBotClient struct {
	client *tbot.Client
}

func (tc *telegramBotClient) SendMessage(chatID string, text string) error {
	_, err := tc.client.SendMessage(chatID, text)
	return err
}

func SendMessageToChat(client BotClientInterface, chatID string, text string) error {
	err := client.SendMessage(chatID, text)
	if err != nil {
		log.Panic(err)
	}
	return err
}

func NewTelegramBotClient(settings settings.Settings) BotClientInterface {
	botSettings := settings.Bot()
	return &telegramBotClient{client: tbot.New(botSettings.Token).Client()}
}
