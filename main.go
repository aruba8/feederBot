package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/biomaks/feederBot/services"
	"github.com/biomaks/feederBot/settings"
	"github.com/biomaks/feederBot/utils"
	"github.com/mmcdole/gofeed"
	"time"
)

func HandleRequest() {
	appSettings := settings.GetSettings()
	feeder := services.NewFeeder(gofeed.NewParser())
	feedParser := utils.NewFeedParser()
	feed, _ := feeder.GetFeed(appSettings.Feeds.Weather)
	alerts := feedParser.ParseFeed(feed)

	mongoClient, _ := services.NewClient(appSettings)
	mongoDb := services.NewDatabase(appSettings, mongoClient)
	mongoStorage := services.NewMongoDatabase(mongoDb)

	checker := utils.NewChecker(mongoStorage)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbAlerts := mongoStorage.FindAllAlerts(ctx, 1, "published", -1)
	checker.Check(alerts, dbAlerts)
}

func main() {
	lambda.Start(HandleRequest)
}
