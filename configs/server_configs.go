package configs

import "os"

const PORT = ":8081"

var CURR_DIR, _ = os.Getwd()

const MONGO_PORT = ":8100"
const MONGO_URI = "mongodb://127.0.0.1" + MONGO_PORT

type redisConfig struct {
	protocol       string
	networkAddress string
	port           string
}

var AuthRedisConfig = redisConfig{
	protocol:       "redis",
	networkAddress: "127.0.0.1",
	port:           "8008",
}

func (rConf redisConfig) GetConnectionURL() string {
	return rConf.protocol + "://" + rConf.networkAddress + ":" + rConf.port
}
