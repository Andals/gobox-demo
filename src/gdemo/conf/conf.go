package conf

import (
	"os"
	"os/user"

	"andals/gobox/exception"
	"andals/gobox/misc"

	"gdemo/errno"
)

var PrjHome string
var Hostname string
var Username string

func Init(prjHomePath string) *exception.Exception {
	if !misc.DirExist(prjHomePath) {
		return exception.New(errno.E_CONF_INVALID_PRJ_HOME, "prjHome not exists")
	}
	PrjHome = prjHomePath

	err := initServerConf(PrjHome + "/conf")
	if err != nil {
		return exception.New(errno.E_CONF_INVALID_SERVER_CONF, "init serverConf error: "+err.Error())
	}

	Hostname, _ = os.Hostname()
	curUser, _ := user.Current()
	Username = curUser.Username

	return nil
}
