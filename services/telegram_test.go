package services

import (
	"github.com/stretchr/testify/mock"
	"github.com/yanzay/tbot/v2"
	"testing"
)

type telegramMock struct {
	mock.Mock
}

func (t *telegramMock) Send(text string, chatID string) (bool, error) {
	args := t.Called(text, chatID)
	return args.Bool(0), args.Error(1)
}

func (t *telegramMock) Client() *tbot.Client {
	args := t.Called()
	return args.Get(0).(*tbot.Client)
}

func (t *telegramMock) Server() *tbot.Server {
	args := t.Called()
	return args.Get(0).(*tbot.Server)
}

func TestTelegramService_SendMessage(t *testing.T) {
	telegramMock := telegramMock{}
	t.Run("test send message", func(t *testing.T) {
		telegramMock.On("Send", "test", "22222222").Return(true, nil)
		telegram := Telegram{&telegramMock}
		telegram.SendMessage("test", "22222222")
		telegramMock.AssertNumberOfCalls(t, "Send", 1)
	})
}
