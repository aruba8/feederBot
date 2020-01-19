package main

import (
	"fmt"
	"github.com/biomaks/feederBot/services"
	"github.com/mmcdole/gofeed"
)

func main() {
	settings := getSettings()
	feeder := services.NewFeeder(gofeed.NewParser())
	fmt.Println(feeder.GetFeed(settings.Feeds.Weather))
	feedParser := NewFeedParser()
	feed, _ := feeder.GetFeed(settings.Feeds.Weather)
	alerts := feedParser.ParseFeed(feed)
	mongoStorage := services.NewMongoStorage("mongodb://localhost:27017", "alerts", "alerts")
	storageService := services.NewStorageService(mongoStorage)

	save(alerts, storageService)

}

func save(alerts []services.Alert, s *services.Storage) {
	for _, alert := range alerts {
		fmt.Println(alert.ID)
		fmt.Println(alert.Title)
		fmt.Println(alert.Datetime)
		fmt.Println(alert.Updated)
		fmt.Println(alert.Published)
		fmt.Println(alert.Categories)
		//_, err := s.Storage.SaveAlert(alert)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
}
