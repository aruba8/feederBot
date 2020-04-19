package utils

import (
	"github.com/biomaks/feederBot/services"
	"log"
)

type Checker struct {
	storageService services.StorageInterface
}

func (c *Checker) Check(feedAlerts []services.Alert, dbAlerts []services.Alert) {
	var err error
	for _, alert := range feedAlerts {
		if len(dbAlerts) > 0 {
			latestDbAlert := dbAlerts[0]
			log.Printf("New alert time: %d. Old alert time: %d", alert.Published.Unix(), latestDbAlert.Published.Unix())
			if alert.Published.Unix() > latestDbAlert.Published.Unix() {
				_, err = c.storageService.SaveAlert(alert)
			}
		} else {
			log.Printf("No previosu alerts. Saving the first alert.")
			_, err = c.storageService.SaveAlert(alert)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewChecker(storageService services.StorageInterface) Checker {
	return Checker{storageService}
}
