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

const collectionName = "alerts"

type MongoClientInterface interface {
	NewClient(*options.ClientOptions) *mongo.Client
}

type CollectionInterface interface {
	Find(ctx context.Context, filter interface{}, opt *options.FindOptions) CursorInterface
	InsertOne(ctx context.Context, document interface{}) (interface{}, error)
}

type DatabaseInterface interface {
	Collection(name string) CollectionInterface
	Client() ClientInterface
}

type ClientInterface interface {
	Database(string) DatabaseInterface
	Connect() error
	StartSession() (mongo.Session, error)
}

type CursorInterface interface {
	Next(ctx context.Context) bool
	Err() error
	Close(ctx context.Context) error
	Decode(interface{}) error
}

type mongoCollection struct {
	coll *mongo.Collection
}
type mongoClient struct {
	client *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoSession struct {
	mongo.Session
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}, opt *options.FindOptions) CursorInterface {
	result, err := mc.coll.Find(ctx, filter, opt)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	result, err := mc.coll.InsertOne(ctx, document)
	return result, err
}

func (c *mongoClient) Database(dbName string) DatabaseInterface {
	d := c.client.Database(dbName)
	return &mongoDatabase{db: d}
}

func (c *mongoClient) Connect() error {
	return c.client.Connect(context.TODO())
}

func (c *mongoClient) StartSession() (mongo.Session, error) {
	s, err := c.client.StartSession()
	return &mongoSession{s}, err
}

func (mdb *mongoDatabase) Client() ClientInterface {
	client := mdb.db.Client()
	return &mongoClient{client}
}

func (mdb *mongoDatabase) Collection(name string) CollectionInterface {
	collection := mdb.db.Collection(name)
	return &mongoCollection{coll: collection}
}

func NewClient(settings settings.Settings, client func(opts ...*options.ClientOptions) (*mongo.Client, error)) (ClientInterface, error) {
	dbSettings := settings.Database()
	c, err := client(options.Client().ApplyURI(dbSettings.ConnectionString()))
	if err != nil {
		log.Panic(err)
	}
	return &mongoClient{client: c}, err
}

func NewDatabase(settings settings.Settings, client ClientInterface) DatabaseInterface {
	return client.Database(settings.Database().DatabaseName)
}

func NewMongoDatabase(db DatabaseInterface) StorageService {
	return &storage{db: db}
}

type StorageService interface {
	SaveAlert(ctx context.Context, alert Alert) (bool, error)
	FindAlerts(ctx context.Context, filter interface{}, sortBy string) []Alert
	FindAllAlerts(ctx context.Context, limit int64, orderBy string, orderDirection int) []Alert
}

type storage struct {
	db DatabaseInterface
}

func (s *storage) SaveAlert(ctx context.Context, alert Alert) (bool, error) {
	_, err := s.db.Collection(collectionName).InsertOne(ctx, alert)
	if err != nil {
		log.Panic(err)
	}
	return true, err
}

func (s *storage) FindAlerts(ctx context.Context, filter interface{}, sortBy string) []Alert {
	opt := options.Find()
	opt.SetSort(bson.D{{
		Key:   sortBy,
		Value: -1,
	}})
	cursor := s.db.Collection(collectionName).Find(ctx, filter, opt)
	defer cursor.Close(ctx)
	var results []Alert
	for cursor.Next(ctx) {
		var alert Alert
		err := cursor.Decode(&alert)
		if err != nil {
			log.Panic(err)
		}
		results = append(results, alert)
	}
	if err := cursor.Err(); err != nil {
		log.Panic(err)
	}
	return results
}

func (s *storage) FindAllAlerts(ctx context.Context, limit int64, orderBy string, orderDirection int) []Alert {
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{
		Key:   orderBy,
		Value: orderDirection,
	}})
	cursor := s.db.Collection(collectionName).Find(ctx, bson.D{}, opts)
	defer cursor.Close(ctx)
	var results []Alert
	for cursor.Next(ctx) {
		var alert Alert
		err := cursor.Decode(&alert)
		if err != nil {
			log.Panic(err)
		}
		results = append(results, alert)
	}
	if err := cursor.Err(); err != nil {
		log.Panic(err)
	}
	return results
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
