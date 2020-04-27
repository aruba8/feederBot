package services

import (
	"errors"
	"github.com/biomaks/feederBot/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockBotClientInterface struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: chatID, text
func (_m *MockBotClientInterface) SendMessage(chatID string, text string) error {
	ret := _m.Called(chatID, text)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(chatID, text)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func TestNewTelegramClient(t *testing.T) {

	var botClient BotClientInterface

	t.Run("test client returning nil", func(t *testing.T) {
		botClient = &MockBotClientInterface{}
		botClient.(*MockBotClientInterface).On("SendMessage", "1111", "tst").
			Return(nil)
		err := SendMessageToChat(botClient, "1111", "tst")
		assert.Nil(t, err)

	})

	t.Run("test client returning error", func(t *testing.T) {
		botClient = &MockBotClientInterface{}
		botClient.(*MockBotClientInterface).On("SendMessage", "1111", "tst").
			Return(errors.New("something went wrong"))
		//err := SendMessageToChat(botClient, "1111", "tst")
		assert.Panics(t, func() {
			_ = SendMessageToChat(botClient, "1111", "tst")
		})

	})

	t.Run("NewTelegramBotClient returns correct structure", func(t *testing.T) {
		conf := settings.GetSettings()
		cl := NewTelegramBotClient(conf)
		assert.IsType(t, &telegramBotClient{}, cl)
	})

}
