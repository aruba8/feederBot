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

func (s *storageMock) FindAlerts(filter interface{}) []*Alert {
	args := s.Called(filter)
	return args.Get(0).([]*Alert)
}

func M(a, b interface{}) []interface{} {
	return [] interface{}{a, b}
}

func TestMongodbStorage(t *testing.T) {

	t.Run("test FindAlerts", func(t *testing.T) {
		theMock := storageMock{}
		theMock.On("FindAlerts", bson.M{}).Return([]*Alert{})

		stService := Storage{&theMock}
		assert.Equal(t, []*Alert{}, stService.Storage.FindAlerts(bson.M{}))
		theMock.AssertNumberOfCalls(t, "FindAlerts", 1)
	})

	t.Run("test SaveAlert", func(t *testing.T) {
		theMock := storageMock{}
		theMock.On("SaveAlert", Alert{}).Return(true, nil)

		stService := Storage{&theMock}
		assert.Equal(t, M(true, nil), M(stService.Storage.SaveAlert(Alert{})))
		theMock.AssertNumberOfCalls(t, "SaveAlert", 1)
	})

}
