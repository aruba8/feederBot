package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/biomaks/feederBot/services"
	"github.com/biomaks/feederBot/settings"
	"github.com/biomaks/feederBot/utils"
	"github.com/mmcdole/gofeed"
)

func HandleRequest() {
	appSettings := settings.GetSettings()
	feeder := services.NewFeeder(gofeed.NewParser())
	feedParser := utils.NewFeedParser()
	feed, _ := feeder.GetFeed(appSettings.Feeds.Weather)
	alerts := feedParser.ParseFeed(feed)

	mongoStorage := services.NewMongoStorage(appSettings)
	storageService := services.NewStorageService(mongoStorage)
	checker := utils.NewChecker(storageService.Storage)
	dbAlerts := storageService.Storage.FindAllAlerts(1, "published", -1)
	checker.Check(alerts, dbAlerts)
}

func main() {
	lambda.Start(HandleRequest)
}
