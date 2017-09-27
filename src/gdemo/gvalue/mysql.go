package gvalue

import (
	"andals/gobox/mysql"

	"gdemo/conf"
)

var MysqlClientPool *mysql.Pool

func InitMysql() {
	MysqlClientPool = mysql.NewPool(conf.MysqlConf.PoolClientTimeout, conf.MysqlConf.PoolSize, NewMysqlClient)
}

func NewMysqlClient() (*mysql.Client, error) {
	config := mysql.NewConfig(conf.MysqlConf.User, conf.MysqlConf.Pass, conf.MysqlConf.Host, conf.MysqlConf.Port, conf.MysqlConf.Name)
	config.LogLevel = conf.LogConf.Level
	config.ReadTimeout = conf.MysqlConf.RWTimeout
	config.WriteTimeout = conf.MysqlConf.RWTimeout

	return mysql.NewClient(config, nil)
}
