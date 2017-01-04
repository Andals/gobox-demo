package conf

import (
	"andals/gobox/exception"

	"gdemo/errno"
)

func Init(prjHome string) *exception.Exception {
	err := initServerConf(prjHome)
	if err != nil {
		return exception.New(errno.E_CONF_INVALID_SERVER_CONF, "init serverConf error: "+err.Error())
	}

	return nil
}
