package utils

import (
	"context"
	"github.com/biomaks/feederBot/services"
	"log"
	"time"
)

type Checker struct {
	storageService services.StorageService
}

func (c *Checker) Check(feedAlerts []services.Alert, dbAlerts []services.Alert) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	for _, alert := range feedAlerts {
		if len(dbAlerts) > 0 {
			latestDbAlert := dbAlerts[0]
			log.Printf("New alert time: %d. Old alert time: %d", alert.Published.Unix(), latestDbAlert.Published.Unix())
			if alert.Published.Unix() > latestDbAlert.Published.Unix() {
				_, err = c.storageService.SaveAlert(ctx, alert)
			}
		} else {
			log.Printf("No previosu alerts. Saving the first alert.")
			_, err = c.storageService.SaveAlert(ctx, alert)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewChecker(storageService services.StorageService) Checker {
	return Checker{storageService}
}
