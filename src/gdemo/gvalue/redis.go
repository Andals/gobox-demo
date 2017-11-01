package gvalue

import (
	"andals/gobox/redis"
	"andals/golog"

	"gdemo/conf"
)

var RedisClientPool *redis.Pool

func InitRedis() {
	RedisClientPool = redis.NewPool(conf.RedisConf.PoolClientTimeout, conf.RedisConf.PoolSize, NewRedisClient)
}

func NewRedisClient() (*redis.Client, error) {
	config := redis.NewConfig(conf.RedisConf.Host, conf.RedisConf.Port, conf.RedisConf.Pass)
	config.LogLevel = golog.LEVEL_DEBUG
	config.ReadTimeout = conf.RedisConf.RWTimeout
	config.WriteTimeout = conf.RedisConf.RWTimeout

	return redis.NewClient(config, nil)
}
