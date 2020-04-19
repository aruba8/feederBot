package utils

import (
	"github.com/biomaks/feederBot/services"
	"log"
)

type Checker struct {
	storageService services.StorageInterface
}

func (c *Checker) Check(feedAlerts []services.Alert, dbAlert services.Alert) {
	for _, alert := range feedAlerts {
		log.Printf("New alert time: %d. Old alert time: %d", alert.Published.Unix(), dbAlert.Published.Unix())
		if alert.Published.Unix() > dbAlert.Published.Unix() {
			_, _ = c.storageService.SaveAlert(alert)
		}
	}
}

func NewChecker(storageService services.StorageInterface) Checker {
	return Checker{storageService}
}
