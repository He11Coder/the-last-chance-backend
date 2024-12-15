package app

import (
	"context"
	"mainService/configs"

	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetMongo() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(configs.MONGO_URI)
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

func GetRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,

		Wait: true,

		IdleTimeout:     0,
		MaxConnLifetime: 0,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(configs.AuthRedisConfig.GetConnectionURL())
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
