package main

import (
	"github.com/biomaks/feederBot/services"
	"github.com/mmcdole/gofeed"
	"log"
)

func main() {
	settings := getSettings()
	feeder := services.NewFeeder(gofeed.NewParser())
	feedParser := NewFeedParser()
	feed, _ := feeder.GetFeed(settings.Feeds.Weather)
	alerts := feedParser.ParseFeed(feed)
	mongoStorage := services.NewMongoStorage(settings.getDbUri(), settings.Database.Name, settings.Database.Collection)
	storageService := services.NewStorageService(mongoStorage)
	check(alerts, storageService)
}

func check(feedAlerts []services.Alert, s *services.Storage) {

	alerts := s.Storage.FindAllAlerts(1, "published", -1)
	if len(alerts) < 1 {
		log.Println("No previous alerts")
		return
	}
	for _, feedAlert := range feedAlerts {
		if feedAlert.Published.Unix() > alerts[0].Published.Unix() {
			s.Storage.SaveAlert(feedAlert)
		} else {
			log.Println("Nothing changed")
		}
	}
}
