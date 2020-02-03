package services

import (
	"github.com/stretchr/testify/mock"
	"github.com/yanzay/tbot"
	"testing"
)

type telegramMock struct {
	mock.Mock
}

func (t *telegramMock) Send(text string, tm tbot.Message) (bool, error) {
	args := t.Called(text, tm)
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
	mock := telegramMock{}
	t.Run("test send message", func(t *testing.T) {
		mock.On("Send", "test", tbot.Message{
			Chat: tbot.Chat{
				ID: "1111",
			},
		}).Return(true, nil)
		telegram := Telegram{&mock}
		telegram.SendMessage("test", tbot.Message{
			Chat: tbot.Chat{
				ID: "1111",
			},
		})
		mock.AssertNumberOfCalls(t, "Send", 1)
	})
}
