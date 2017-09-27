package conf

import "time"

var MysqlConf struct {
	Host              string
	User              string
	Pass              string
	Port              string
	Name              string
	RWTimeout         time.Duration
	PoolClientTimeout time.Duration
	PoolSize          int
}

func initMysqlConf() {
	MysqlConf.Host = scJson.Mysql.Host
	MysqlConf.User = scJson.Mysql.User
	MysqlConf.Pass = scJson.Mysql.Pass
	MysqlConf.Port = scJson.Mysql.Port
	MysqlConf.Name = scJson.Mysql.Name
	MysqlConf.RWTimeout = time.Duration(scJson.Mysql.RWTimeoutSeconds) * time.Second
	MysqlConf.PoolClientTimeout = time.Duration(scJson.Mysql.PoolClientTimeoutSeconds) * time.Second
	MysqlConf.PoolSize = scJson.Mysql.PoolSize
}
