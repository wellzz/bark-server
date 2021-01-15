package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type DB struct {
	db       *mongo.Client
	cache    map[string]string
	alias    map[string]string
	database string
	table    string
}

const (
	// environment variables
	mongoDBConnectionStringEnvVarName = "MONGODB_CONNECTION_STRING"
	mongoDBDatabaseEnvVarName         = "MONGODB_DATABASE"
	mongoDBCollectionEnvVarName       = "MONGODB_COLLECTION"
)

func NewDB() (*DB, error) {
	conn := os.Getenv(mongoDBConnectionStringEnvVarName)
	database := os.Getenv(mongoDBDatabaseEnvVarName)
	collection := os.Getenv(mongoDBCollectionEnvVarName)
	if database == "" {
		database = "base"
	}
	if collection == "" {
		collection = "bark"
	}
	return newDB(conn, database, collection)
}

func newDB(conn, database, collection string) (*DB, error) {
	if conn == "" {
		return nil, errors.New("missing conn")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	clientOptions := options.Client().ApplyURI(conn).SetDirect(true)
	c, err := mongo.NewClient(clientOptions)

	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		db:       c,
		cache:    make(map[string]string),
		alias:    make(map[string]string),
		database: database,
		table:    collection,
	}, nil
}

func (d *DB) Get(key string) (string, error) {
	if v, ok := d.cache[key]; ok {
		return v, nil
	}
	device, err := d.get(key)
	if err != nil {
		return "", err
	}
	d.cache[key] = device.Token
	return device.Token, nil
}

func (d *DB) Set(key, token string) (string, error) {
	err := d.set(key, token, d.alias[key])
	if err != nil {
		return "", err
	}
	if k := d.alias[key]; len(k) > 0 {
		d.cache[k] = token
		delete(d.cache, key)
		delete(d.alias, key)
		return k, nil
	} else {
		d.cache[key] = token
	}
	return key, nil
}

func (d *DB) GetAlias(key string) string {
	return d.alias[key]
}

func (d *DB) Alias(key, alias string) error {
	if d, err := d.get(alias); err != nil && err != mongo.ErrNoDocuments {
		return err
	} else if d != nil {
		return errors.New("duplicate key")
	}
	for k, v := range d.alias {
		if v == alias && k != key {
			return errors.New("duplicate key")
		}
	}
	d.alias[key] = alias
	return nil
}

func (d *DB) Close() error {
	return d.db.Disconnect(context.Background())
}

func (d *DB) get(key string) (*Device, error) {
	filter := bson.D{{"key", key}}
	rs := d.collection().FindOne(context.Background(), filter)
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	t := &Device{}
	err := rs.Decode(t)
	return t, err
}

func (d *DB) set(key, token, alias string) error {
	filter := bson.D{{"key", key}}
	newKey := key
	if len(alias) != 0 {
		newKey = alias
	}
	update := bson.D{
		{"$set", bson.D{
			{"key", newKey},
			{"token", token},
		}},
	}
	opts := options.Update().SetUpsert(true)
	_, err := d.collection().UpdateOne(context.Background(), filter, update, opts)
	return err
}

func (d *DB) collection() *mongo.Collection {
	return d.db.Database(d.database).Collection(d.table)
}

type Device struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Key   string             `bson:"key"`
	Token string             `bson:"token"`
}
