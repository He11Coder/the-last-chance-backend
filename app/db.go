package app

import (
	"context"
	"mainService/configs"

	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetMongo() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(configs.MainMongoConfig.GetConnectionURI())
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

func InitDBAndIndexes(cli *mongo.Client) (*mongo.Database, error) {
	db := cli.Database("tlc")

	serviceColl := db.Collection("service")
	textIndex := mongo.IndexModel{
		Keys: bson.D{
			{"title", "text"},
			{"description", "text"},
		},
		Options: options.Index().
			SetName("textIndex").
			SetWeights(bson.M{
				"title":       10,
				"description": 5,
			}).
			SetDefaultLanguage("russian"),
	}

	_, err := serviceColl.Indexes().CreateOne(context.TODO(), textIndex)
	if err != nil {
		return nil, err
	}

	userColl := db.Collection("user")
	loginIndex := mongo.IndexModel{
		Keys: bson.D{
			{"login", 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetName("loginIndex"),
	}

	_, err = userColl.Indexes().CreateOne(context.TODO(), loginIndex)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,

		Wait: true,

		IdleTimeout:     0,
		MaxConnLifetime: 0,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(configs.AuthRedisConfig.GetConnectionURI())
			if err != nil {
				return nil, err
			}

			_, err = redis.String(conn.Do("PING"))
			if err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
	}

	return pool
}
