package conf

import (
	"encoding/json"
	"io/ioutil"

	"andals/gobox/exception"
	"andals/gobox/misc"

	"gdemo/errno"
)

var ServerConf serverConf

type serverConf struct {
	PrjName string
	IsDev   bool

	FrontDomain     string
	FrontGoHttpPort string

	DataRoot string
	LogRoot  string
	TmpRoot  string

	FrontPidFile string
}

func initServerConf(confRoot string) *exception.Exception {
	scJson := new(serverConfJson)

	err := parseServerConfJson(confRoot+"/server/server_conf.json", scJson)
	if err != nil {
		return err
	}
	err = parseServerConfJson(confRoot+"/server_conf_rewrite.json", scJson)
	if err != nil {
		return err
	}

	ServerConf.PrjName = scJson.PrjName
	ServerConf.IsDev = scJson.IsDev

	ServerConf.FrontDomain = scJson.FrontDomain
	ServerConf.FrontGoHttpPort = scJson.FrontGoHttpPort

	ServerConf.DataRoot = PrjHome + "/data"
	ServerConf.LogRoot = PrjHome + "/logs"
	ServerConf.TmpRoot = PrjHome + "/tmp"

	ServerConf.FrontPidFile = ServerConf.TmpRoot + "/front.pid"

	return nil
}

type serverConfJson struct {
	PrjName string `json:"prj_name"`
	IsDev   bool   `json:"is_dev"`

	FrontDomain     string `json:"front_domain"`
	FrontGoHttpPort string `json:"front_gohttp_port"`
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
