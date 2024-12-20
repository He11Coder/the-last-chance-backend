package configs

import (
	"os"
)

var PORT = ":"

var CURR_DIR, _ = os.Getwd()

type dbConfig struct {
	protocol string
	host     string
	port     string
}

var MainMongoConfig = dbConfig{}

var AuthRedisConfig = dbConfig{}

func InitConfigs() {
	PORT = PORT + os.Getenv("MAIN_SERVICE_PORT")

	MainMongoConfig.protocol = os.Getenv("MONGO_PROTOCOL")
	MainMongoConfig.host = os.Getenv("MONGO_HOST")
	MainMongoConfig.port = os.Getenv("MONGO_PORT")

	AuthRedisConfig.protocol = os.Getenv("REDIS_PROTOCOL")
	AuthRedisConfig.host = os.Getenv("REDIS_HOST")
	AuthRedisConfig.port = os.Getenv("REDIS_PORT")
}

func (conf dbConfig) GetConnectionURI() string {
	return conf.protocol + "://" + conf.host + ":" + conf.port
}
