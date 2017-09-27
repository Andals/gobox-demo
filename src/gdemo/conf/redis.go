package conf

import "time"

var RedisConf struct {
	Host              string
	Pass              string
	Port              string
	RWTimeout         time.Duration
	PoolClientTimeout time.Duration
	PoolSize          int
}

func initRedisConf() {
	RedisConf.Host = scJson.Redis.Host
	RedisConf.Pass = scJson.Redis.Pass
	RedisConf.Port = scJson.Redis.Port
	RedisConf.RWTimeout = time.Duration(scJson.Redis.RWTimeoutSeconds) * time.Second
	RedisConf.PoolClientTimeout = time.Duration(scJson.Redis.PoolClientTimeoutSeconds) * time.Second
	RedisConf.PoolSize = scJson.Redis.PoolSize
}
