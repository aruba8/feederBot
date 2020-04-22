package services

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"net/url"
)

type FeedService interface {
	GetFeed(url string) (*gofeed.Feed, error)
}

type ParserInterface interface {
	ParseURL(url string) (*gofeed.Feed, error)
}

type Feeder struct {
	Parser ParserInterface
}

func (f *Feeder) GetFeed(url string) (*gofeed.Feed, error) {
	if !isValidUrl(url) {
		return nil, errors.New("not valid URL")
	}

	feed, err := f.Parser.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return feed, err
}

func NewFeeder(p ParserInterface) FeedService {
	return &Feeder{p}
}

func isValidUrl(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	} else {
		return true
	}
}
