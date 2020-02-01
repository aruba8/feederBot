package main

import (
	"context"
	"github.com/biomaks/feederBot/services"
	"github.com/biomaks/feederBot/settings"
	"github.com/mmcdole/gofeed"
	"log"
)

func HandleRequest(ctx context.Context) {
	settings := settings.GetSettings()
	feeder := services.NewFeeder(gofeed.NewParser())
	feedParser := NewFeedParser()
	feed, _ := feeder.GetFeed(settings.Feeds.Weather)
	alerts := feedParser.ParseFeed(feed)
	mongoStorage := services.NewMongoStorage(settings)
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

func main() {
	//lambda.Start(HandleRequest)
	settings := settings.GetSettings()
	feeder := services.NewFeeder(gofeed.NewParser())
	feedParser := NewFeedParser()
	feed, _ := feeder.GetFeed(settings.Feeds.Weather)
	alerts := feedParser.ParseFeed(feed)
	mongoStorage := services.NewMongoStorage(settings)
	s := services.NewStorageService(mongoStorage)
	//for _, v := range alerts{
	//	println(v.Description)
	//	als := s.Storage.FindAllAlerts(1, "published", -1)
	//	println(len(als))
	//}
	check(alerts, s)

}
