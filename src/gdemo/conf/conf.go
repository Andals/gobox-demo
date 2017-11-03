package conf

import (
	"github.com/andals/gobox/exception"
	gmisc "github.com/andals/gobox/misc"

	"gdemo/errno"
)

var PrjHome string

func Init(prjHome string) *exception.Exception {
	if !gmisc.DirExist(prjHome) {
		return exception.New(errno.E_SYS_INVALID_PRJ_HOME, "prjHome not exists")
	}

	PrjHome = prjHome

	err := initServerConfJson()
	if err != nil {
		return exception.New(errno.E_SYS_INIT_SERVER_CONF_JSON_FAIL, "init serverConfJson error: "+err.Error())
	}

	initBaseConf()
	initLogConf()
	initPprofConf()
	initHttpConf()
	initRedisConf()
	initMysqlConf()

	return nil
}
