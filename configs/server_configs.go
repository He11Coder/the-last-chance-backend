package configs

import "os"

const PORT = ":8081"

var CURR_DIR, _ = os.Getwd()

type redisConfig struct {
	protocol       string
	networkAddress string
	port           string
}

var AuthRedisConfig = redisConfig{
	protocol:       "redis",
	networkAddress: "sessions_hnh",
	port:           "6379",
}

func (rConf redisConfig) GetConnectionURL() string {
	return rConf.protocol + "://" + "@" + rConf.networkAddress + ":" + rConf.port
}
