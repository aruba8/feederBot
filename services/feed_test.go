package services

import (
	"errors"
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

type MockParserInterface struct {
	mock.Mock
}

// ParseURL provides a mock function with given fields: url
func (_m *MockParserInterface) ParseURL(url string) (*gofeed.Feed, error) {
	ret := _m.Called(url)

	var r0 *gofeed.Feed
	if rf, ok := ret.Get(0).(func(string) *gofeed.Feed); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gofeed.Feed)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func TestFeeder(t *testing.T) {
	t.Run("when Parser returns error", func(t *testing.T) {

		var parserMock ParserInterface

		parserMock = &MockParserInterface{}

		feed := &gofeed.Feed{}
		parserMock.(*MockParserInterface).
			On("ParseURL", mock.Anything).
			Return(feed, errors.New("something happen"))

		feeder := NewFeeder(parserMock)
		assert.Panics(t, func() {feeder.GetFeed("http://google.com")}, "")



	})

	t.Run("when url is not valid", func(t *testing.T) {

		notValidUrls := []string{
			"", "http", "ss", "http/://google.com",
		}
		theParserMock := parserMock{}
		f := Feeder{&theParserMock}
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
			f := Feeder{&theParserMock}
			_, err := f.GetFeed(url)
			assert.Equal(t, nil, err)
			theParserMock.AssertNumberOfCalls(t, "ParseURL", 1)

		}
	})
}
