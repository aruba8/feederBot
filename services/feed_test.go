package services

import (
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type parserMock struct {
	mock.Mock
}

func (f *parserMock) ParseURL(url string) (*gofeed.Feed, error) {
	args := f.Called(url)
	return args.Get(0).(*gofeed.Feed), args.Error(1)
}

func TestFeeder(t *testing.T) {
	t.Run("when url is not valid", func(t *testing.T) {

		notValidUrls := [] string{
			"", "http", "ss", "http/://google.com",
		}
		theParserMock := parserMock{}
		f := feeder{&theParserMock}
		for _, url := range notValidUrls {
			_, err := f.GetFeed(url)
			assert.Equal(t, "not valid URL", err.Error())
			theParserMock.AssertNumberOfCalls(t, "ParseURL", 0)
		}
	})

	t.Run("when url is valid", func(t *testing.T) {

		validUrls := []string{
			"http://test.com", "https://google.com/xml",
		}
		for _, url := range validUrls {
			theParserMock := parserMock{}
			theParserMock.On("ParseURL", url).Return(&gofeed.Feed{}, nil)
			f := feeder{&theParserMock}
			_, err := f.GetFeed(url)
			assert.Equal(t, nil, err)
			theParserMock.AssertNumberOfCalls(t, "ParseURL", 1)

		}
	})
}
