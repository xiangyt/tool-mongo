package output

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const MongoUri = "mongodb://localhost:27017"

// MongoDB 数据存储
type Store struct {
	dbname string
	client *mongo.Client
}

func (s *Store) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoUri))
	if err != nil {
		log.Printf("Connect mongo failed! err:%s\n", err)
		return err
	} else {
		s.client = client
	}

	//defer client.Disconnect()
	// 检查连接
	err = s.client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Ping mongo failed! err:%s\n", err)
		return err
	}
	log.Println("Connected to MongoDB!")

	return nil
}

func (s *Store) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Disconnect(context.Background())
}

func (s *Store) Client() *mongo.Client {
	return s.client
}

func (s *Store) DataBase(name ...string) *mongo.Database {
	if len(name) > 0 && name[0] != "" {
		return s.client.Database(name[0])
	}
	return s.client.Database(s.dbname)
}

func (s *Store) Collection(name string, dbName ...string) *mongo.Collection {
	var db = s.dbname
	if len(dbName) > 0 && dbName[0] != "" {
		db = dbName[0]
	}
	return s.client.Database(db).Collection(name)
}

const DefaultDBName = "share"

var s *Store

func NewStore() *Store {
	s = &Store{
		dbname: DefaultDBName,
		client: nil,
	}
	return s
}

func GetCollection(name string) *mongo.Collection {
	return s.Collection(name)
}
