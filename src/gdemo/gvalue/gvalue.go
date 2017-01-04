package gvalue

import (
	"andals/gobox/exception"

	"family/errno"
)

func Init() *exception.Exception {
	err := initHuajiao()
	if err != nil {
		return exception.New(errno.E_SYS_INIT_HUAJIAO_FAIL, err.Error())
	}

	err = initAthena()
	if err != nil {
		return exception.New(errno.E_SYS_INIT_ATHENAPAGE_FAIL, err.Error())
	}

	go RunHuajiaoLiveDataTask()
	go RunAthenaPageTask()

	return nil
}
