package services

import (
	"context"
	"errors"
	"github.com/biomaks/feederBot/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {

	var dbMock DatabaseInterface
	var collectionMock CollectionInterface
	var cursorMock CursorInterface

	dbMock = &MockDatabaseInterface{}
	collectionMock = &MockCollectionInterface{}
	cursorMock = &MockCursorInterface{}

	t.Run("test FindAlerts", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		collectionMock.(*MockCollectionInterface).
			On("Find", ctx, bson.M{}, options.Find().SetSort(bson.D{{
				Key:   "",
				Value: -1,
			}})).
			Return(cursorMock)
		dbMock.(*MockDatabaseInterface).
			On("Collection", "alerts").
			Return(collectionMock)

		cursorMock.(*MockCursorInterface).On("Next", ctx).Return(false)
		cursorMock.(*MockCursorInterface).On("Close", ctx).Return(nil)
		cursorMock.(*MockCursorInterface).On("Err").Return(nil)

		alertDb := NewMongoDatabase(dbMock)

		alerts := alertDb.FindAlerts(ctx, bson.M{}, "")
		assert.Empty(t, alerts)
	})

	t.Run("test FindAlerts error", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		cursorMock.(*MockCursorInterface).On("Next", ctx).Return(false)
		cursorMock.(*MockCursorInterface).On("Close", ctx).Return(nil)
		cursorMock.(*MockCursorInterface).On("Err").Return(errors.New("failed"))
		collectionMock.(*MockCollectionInterface).
			On("Find", ctx, bson.D{}, options.Find().SetSort(bson.D{{
				Key:   "",
				Value: -1,
			}})).
			Return(cursorMock)
		dbMock.(*MockDatabaseInterface).
			On("Collection", "alerts").
			Return(collectionMock)

		alertDb := NewMongoDatabase(dbMock)

		assert.Panics(t, func() {
			alertDb.FindAlerts(ctx, bson.M{}, "")
		})
	})

	t.Run("test SaveAlert", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		alert := Alert{}

		collectionMock.(*MockCollectionInterface).
			On("InsertOne", ctx, alert).
			Return(bson.M{}, nil)
		dbMock.(*MockDatabaseInterface).
			On("Collection", "alerts").
			Return(collectionMock)
		alertDb := NewMongoDatabase(dbMock)
		result, err := alertDb.SaveAlert(ctx, alert)

		assert.Nil(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("test SaveAlert panics", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		alert := Alert{}

		collectionMock.(*MockCollectionInterface).
			On("InsertOne", ctx, alert).
			Return(nil, errors.New("error on test SaveAlert panics"))
		dbMock.(*MockDatabaseInterface).
			On("Collection", "alerts").
			Return(collectionMock)
		alertDb := NewMongoDatabase(dbMock)

		assert.Panics(t, func() {
			alertDb.SaveAlert(ctx, alert)
		})

	})

	t.Run("test FindAllAlerts", func(t *testing.T) {

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		option := options.Find()
		option.SetLimit(1)
		option.SetSort(bson.D{{
			Key:   "",
			Value: 1,
		}})

		collectionMock.(*MockCollectionInterface).
			On("Find", ctx, bson.D{}, option).
			Return(cursorMock)
		dbMock.(*MockDatabaseInterface).
			On("Collection", "alerts").
			Return(collectionMock)

		cursorMock.(*MockCursorInterface).On("Next", ctx).Return(false)
		cursorMock.(*MockCursorInterface).On("Close", ctx).Return(nil)
		cursorMock.(*MockCursorInterface).On("Err").Return(nil)

		alertDb := NewMongoDatabase(dbMock)

		alerts := alertDb.FindAllAlerts(ctx, 1, "", 1)
		assert.Empty(t, alerts)
	})

	t.Run("test NewDatabase", func(t *testing.T) {
		os.Setenv("SETTINGS_FILE_PATH", "../settings.toml")
		var clientInterface ClientInterface
		clientInterface = &MockClientInterface{}
		dbMock = &MockDatabaseInterface{}

		testSettings := settings.GetSettings()

		clientInterface.(*MockClientInterface).
			On("Database", mock.Anything).
			Return(dbMock)

		db := NewDatabase(testSettings, clientInterface)

		assert.Empty(t, db)
		assert.IsType(t, &MockDatabaseInterface{}, db)

	})

	t.Run("test NewClient panics", func(t *testing.T) {
		clientFunc := func(opts ...*options.ClientOptions) (*mongo.Client, error) {
			return nil, errors.New("error")
		}
		os.Setenv("SETTINGS_FILE_PATH", "../settings.toml")
		var clientInterface ClientInterface
		clientInterface = &MockClientInterface{}
		dbMock = &MockDatabaseInterface{}
		testSettings := settings.GetSettings()
		clientInterface.(*MockClientInterface).
			On("Database", mock.Anything).
			Return(dbMock)
		assert.Panics(t, func() {
			NewClient(testSettings, clientFunc)
		})

	})

	t.Run("test NewClient returns client", func(t *testing.T) {
		clientFunc := func(opts ...*options.ClientOptions) (*mongo.Client, error) {
			return &mongo.Client{}, nil
		}
		os.Setenv("SETTINGS_FILE_PATH", "../settings.toml")
		var clientInterface ClientInterface
		clientInterface = &MockClientInterface{}
		dbMock = &MockDatabaseInterface{}
		testSettings := settings.GetSettings()
		clientInterface.(*MockClientInterface).
			On("Database", mock.Anything).
			Return(dbMock)
		client, err := NewClient(testSettings, clientFunc)
		assert.Nil(t, err)
		assert.IsType(t, &mongoClient{}, client)
	})

}
