package utils

import (
	"github.com/biomaks/feederBot/services"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
	"time"
)

type storageMock struct {
	mock.Mock
}

func (s *storageMock) SaveAlert(alert services.Alert) (bool, error) {
	args := s.Called(alert)
	return args.Bool(0), args.Error(1)
}

func (s *storageMock) FindAlerts(filter interface{}, sortBy string) []services.Alert {
	args := s.Called(filter)
	return args.Get(0).([]services.Alert)
}

func (s *storageMock) FindAllAlerts(limit int64, orderBy string, orderDirection int) []services.Alert {
	args := s.Called(limit, orderBy, orderDirection)
	return args.Get(0).([]services.Alert)
}

func createFeedAlerts() []services.Alert {
	var alerts []services.Alert
	for i := 0; i < 10; i++ {
		alerts = append(alerts, createAlert(time.Now().Add(time.Duration(-i)*time.Hour)))
	}
	return alerts
}

func createAlert(publishedTime time.Time) services.Alert {
	alert := services.Alert{
		ID:          primitive.ObjectID{},
		EntryId:     "",
		FeedLink:    "",
		Title:       "",
		Datetime:    time.Now(),
		Updated:     time.Now(),
		Published:   publishedTime,
		Categories:  nil,
		Description: "",
	}
	return alert
}

func TestChecker(t *testing.T) {
	storageServiceMock := storageMock{}
	t.Run("test checker when feed alert is newer than db alert", func(t *testing.T) {
		storageServiceMock.On("SaveAlert", mock.Anything).Return(true, nil)
		checker := NewChecker(&storageServiceMock)
		feedAlerts := createFeedAlerts()
		dbAlerts := []services.Alert{createAlert(time.Now().Add(time.Duration(-24) * time.Hour))}
		log.Println("test checker when feed alert is newer than db alert")
		checker.Check(feedAlerts, dbAlerts)
		storageServiceMock.AssertNumberOfCalls(t, "SaveAlert", 10)
	})
}

func TestNewChecker(t *testing.T) {
	storageServiceMock := storageMock{}
	t.Run("test checker when feed alert is older than db alert", func(t *testing.T) {
		storageServiceMock.On("SaveAlert", mock.Anything).Return(true, nil)
		checker := NewChecker(&storageServiceMock)
		feedAlerts := createFeedAlerts()
		dbAlerts := []services.Alert{createAlert(time.Now())}
		log.Println("test checker when feed alert is older than db alert")
		checker.Check(feedAlerts, dbAlerts)
		storageServiceMock.AssertNumberOfCalls(t, "SaveAlert", 0)
	})
}
