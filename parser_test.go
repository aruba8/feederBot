package main

import (
	"fmt"
	"github.com/biomaks/feederBot/services"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestParser_ParseFeed(t *testing.T) {
	feed := gofeed.Feed{
		Title:           "Test title",
		Description:     "test description",
		Link:            "link",
		FeedLink:        "feed link",
		Updated:         "Updated",
		UpdatedParsed:   nil,
		Published:       "Published",
		PublishedParsed: nil,
		Author:          nil,
		Language:        "Language",
		Image:           nil,
		Copyright:       "Copyright",
		Generator:       "Generator",
		Categories:      nil,
		DublinCoreExt:   nil,
		ITunesExt:       nil,
		Extensions:      nil,
		Custom:          nil,
		Items:           nil,
		FeedType:        "FeedType",
		FeedVersion:     "FeedVersion",
	}

	t.Run("test parser empty items", func(t *testing.T) {
		parser := NewFeedParser()
		got := parser.ParseFeed(&feed)
		assert.Equal(t, 0, len(got))
	})

	t.Run("test parser not empty items", func(t *testing.T) {
		updatedTime := time.Now()
		item := gofeed.Item{
			Title:           "Title",
			Description:     "test description",
			Content:         "Content",
			Link:            "Link",
			Updated:         updatedTime.String(),
			UpdatedParsed:   nil,
			Published:       "Published",
			PublishedParsed: nil,
			Author:          nil,
			GUID:            "",
			Image:           nil,
			Categories:      nil,
			Enclosures:      nil,
			DublinCoreExt:   nil,
			ITunesExt:       nil,
			Extensions:      nil,
			Custom:          nil,
		}
		var items []*gofeed.Item
		items = append(items, &item)
		feed.Items = items

		var want []services.Alert

		alert := services.Alert{
			ID:          primitive.ObjectID{},
			FeedLink:    "feed link",
			Title:       "Title",
			Datetime:    time.Time{},
			Updated:     updatedTime,
			Published:   time.Time{},
			Categories:  nil,
			Description: "test description",
		}
		want = append(want, alert)
		parser := NewFeedParser()
		got := parser.ParseFeed(&feed)
		fmt.Println(got)
		assert.Equal(t, 1, len(got))
		assert.Equal(t, want[0].Published, got[0].Published)
		assert.Equal(t, want[0].EntryId, got[0].EntryId)
		assert.Equal(t, want[0].Categories, got[0].Categories)
		//assert.Equal(t, want[0].Updated, got[0].Updated)
		//assert.Equal(t, want[0].Datetime, got[0].Datetime)
		assert.Equal(t, want[0].Title, got[0].Title)
		assert.Equal(t, want[0].FeedLink, got[0].FeedLink)
		assert.Equal(t, want[0].ID, got[0].ID)
		assert.Equal(t, want[0].Description, got[0].Description)

	})
}
