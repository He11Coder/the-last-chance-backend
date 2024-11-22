package app

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const MONGO_PORT = ":27017"
const MONGO_URI = "mongodb://127.0.0.1" + MONGO_PORT

func GetMongo() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
