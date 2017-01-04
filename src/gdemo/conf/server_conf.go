package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"

	"andals/gobox/exception"
	"andals/gobox/misc"

	"gdemo/errno"
)

var ServerConf serverConf

type serverConf struct {
	PrjHome  string
	Hostname string
	Username string

	PrjName string
	IsDev   bool

	ApiDomain     string
	ApiGoHttpHost string
	ApiGoHttpPort string

	DataRoot string
	LogRoot  string
	TmpRoot  string

	ApiPidFile string
}

func initServerConf(prjHome string) *exception.Exception {
	if !misc.DirExist(prjHome) {
		return exception.New(errno.E_CONF_INVALID_PRJ_HOME, "prjHome not exists")
	}

	confRoot := prjHome + "/conf"
	scJson := new(serverConfJson)

	err := parseServerConfJson(confRoot+"/server/server_conf.json", scJson)
	if err != nil {
		return err
	}
	err = parseServerConfJson(confRoot+"/server_conf_rewrite.json", scJson)
	if err != nil {
		return err
	}

	ServerConf.PrjHome = prjHome
	ServerConf.Hostname, _ = os.Hostname()
	curUser, _ := user.Current()
	ServerConf.Username = curUser.Username

	ServerConf.PrjName = scJson.PrjName
	ServerConf.IsDev = scJson.IsDev

	ServerConf.ApiDomain = scJson.ApiDomain
	ServerConf.ApiGoHttpHost = scJson.ApiGoHttpHost
	ServerConf.ApiGoHttpPort = scJson.ApiGoHttpPort

	ServerConf.DataRoot = ServerConf.PrjHome + "/data"
	ServerConf.LogRoot = ServerConf.PrjHome + "/logs"
	ServerConf.TmpRoot = ServerConf.PrjHome + "/tmp"

	ServerConf.ApiPidFile = ServerConf.TmpRoot + "/api.pid"

	return nil
}

type serverConfJson struct {
	PrjName string `json:"prj_name"`
	IsDev   bool   `json:"is_dev"`

	ApiDomain     string `json:"api_domain"`
	ApiGoHttpHost string `json:"api_gohttp_host"`
	ApiGoHttpPort string `json:"api_gohttp_port"`
}

func parseServerConfJson(path string, scJson *serverConfJson) *exception.Exception {
	if !misc.FileExist(path) {
		return exception.New(errno.E_CONF_INVALID_SERVER_CONF, "file "+path+" not exists")
	}

	jsonStr, err := ioutil.ReadFile(path)
	if err != nil {
		return exception.New(errno.E_CONF_INVALID_SERVER_CONF, err.Error())
	}

	err = json.Unmarshal(jsonStr, scJson)
	if err != nil {
		return exception.New(errno.E_CONF_INVALID_SERVER_CONF, err.Error())
	}

	return nil
}
