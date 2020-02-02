package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type storageMock struct {
	mock.Mock
}

func (s *storageMock) SaveAlert(alert Alert) (bool, error) {
	args := s.Called(alert)
	return args.Bool(0), args.Error(1)
}

func (s *storageMock) FindAlerts(filter interface{}, sortBy string) []Alert {
	args := s.Called(filter)
	return args.Get(0).([]Alert)
}

func (s *storageMock) FindAllAlerts(limit int64, orderBy string, orderDirection int) []Alert {
	args := s.Called(limit, orderBy, orderDirection)
	return args.Get(0).([]Alert)
}

func M(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}

func TestMongodbStorage_FindAlerts(t *testing.T) {

	t.Run("test FindAlerts", func(t *testing.T) {
		theMock := storageMock{}
		theMock.On("FindAlerts", bson.M{}).Return([]Alert{})

		stService := Storage{&theMock}
		assert.Equal(t, []Alert{}, stService.Storage.FindAlerts(bson.M{}, ""))
		theMock.AssertNumberOfCalls(t, "FindAlerts", 1)
		theMock.AssertExpectations(t)
	})
}

func TestMongodbStorage_SaveAlert(t *testing.T) {
	t.Run("test SaveAlert", func(t *testing.T) {
		theMock := storageMock{}
		theMock.On("SaveAlert", Alert{}).Return(true, nil)
		stService := Storage{&theMock}
		assert.Equal(t, M(true, nil), M(stService.Storage.SaveAlert(Alert{})))
		theMock.AssertNumberOfCalls(t, "SaveAlert", 1)
		theMock.AssertExpectations(t)
	})
}

func TestMongodbStorage_FindAllAlerts(t *testing.T) {
	t.Run("test FindAllAlerts", func(t *testing.T) {
		theMock := storageMock{}
		theMock.On("FindAllAlerts", int64(1), "", 1).Return([]Alert{})
		stService := Storage{&theMock}
		assert.Equal(t, []Alert{}, stService.Storage.FindAllAlerts(int64(1), "", 1))
		theMock.AssertNumberOfCalls(t, "FindAllAlerts", 1)
		theMock.AssertExpectations(t)
	})
}
