package services

import (
	"context"
	"github.com/biomaks/feederBot/settings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type StorageInterface interface {
	SaveAlert(alert Alert) (bool, error)
	FindAlerts(interface{}, string) []Alert
	FindAllAlerts(limit int64, orderBy string, orderDirection int) []Alert
}

type MongodbStorage struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Storage struct {
	Storage StorageInterface
}

func (m *MongodbStorage) SaveAlert(alert Alert) (bool, error) {
	_, err := m.collection.InsertOne(m.ctx, alert)
	if err != nil {
		return false, err
	} else {
		return true, err
	}
}

func (m *MongodbStorage) FindAlerts(filter interface{}, sortBy string) []Alert {
	opts := options.Find()
	opts.SetSort(bson.D{{
		Key:   sortBy,
		Value: -1,
	}})

	cursor, err := m.collection.Find(m.ctx, filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(m.ctx)
	var results []Alert
	for cursor.Next(m.ctx) {
		var alert Alert
		err := cursor.Decode(&alert)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, alert)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return results
}

func (m *MongodbStorage) FindAllAlerts(limit int64, orderBy string, orderDirection int) []Alert {
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{
		Key:   orderBy,
		Value: orderDirection,
	}})
	cursor, err := m.collection.Find(m.ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(m.ctx)
	var results []Alert
	for cursor.Next(m.ctx) {
		var alert Alert
		err := cursor.Decode(&alert)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, alert)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return results
}

func NewMongoStorage(settings settings.Settings) StorageInterface {
	mongoDb := settings.Database()
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDb.ConnectionString()))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(settings.Database().Name)
	collection := db.Collection(settings.Database().Name)
	return &MongodbStorage{collection, ctx}
}

func NewStorageService(storageImpl StorageInterface) *Storage {
	return &Storage{storageImpl}
}

type Alert struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EntryId     string             `json:"entry_id" bson:"entry_id"`
	FeedLink    string             `json:"feed_link,omitempty" bson:"feed_link,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	Datetime    time.Time          `json:"datetime,omitempty" bson:"datetime,omitempty"`
	Updated     time.Time          `json:"updated,omitempty" bson:"updated,omitempty"`
	Published   time.Time          `json:"published,omitempty" bson:"published,omitempty"`
	Categories  []string           `json:"categories" bson:"categories,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
}
