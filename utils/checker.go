package utils

import "github.com/biomaks/feederBot/services"

type Checker struct {
	storageService services.StorageInterface
}

func (c *Checker) Check(feedAlerts []services.Alert, dbAlert services.Alert) {
	for _, alert := range feedAlerts {
		if alert.Published.Unix() > dbAlert.Published.Unix() {
			_, _ = c.storageService.SaveAlert(alert)
		}
	}
}

func NewChecker(storageService services.StorageInterface) Checker {
	return Checker{storageService}
}
