package main

import (
	"github.com/biomaks/feederBot/services"
	"github.com/mmcdole/gofeed"
	"time"
)

type FeedParser interface {
	ParseFeed(feed *gofeed.Feed) []services.Alert
}

type Parser struct {
}

func (p *Parser) ParseFeed(feed *gofeed.Feed) []services.Alert {
	items := feed.Items
	var results []services.Alert
	for _, item := range items {

		alert := convert(item, feed)
		results = append(results, alert)
	}
	return results
}

func convert(item *gofeed.Item, feed *gofeed.Feed) services.Alert {
	var alert services.Alert
	updated, _ := time.Parse(time.RFC3339, item.Updated)
	published, _ := time.Parse(time.RFC3339, item.Published)

	alert.FeedLink = feed.FeedLink
	alert.Title = item.Title
	alert.Datetime = time.Now()
	alert.Updated = updated
	alert.Published = published
	alert.Categories = item.Categories

	return alert
}

func NewFeedParser() FeedParser {
	return &Parser{}
}
