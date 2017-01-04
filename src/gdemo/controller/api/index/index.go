package index

import (
	"andals/gobox/http/controller"
	glog "andals/gobox/log"

	"gdemo/misc"
)

func RegAction(cl *controller.Controller) {
	cl.AddExactMatchAction("/api/index", Index)
}

func Index(context *controller.Context, args []string) {
	l, ok := context.TransData[misc.TRANS_DATA_KEY_USER_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Info(context.Rid)
		}
	}

	var data struct {
		Rid []byte
	}
	data.Rid = context.Rid

	context.RespBody = misc.ApiJson(data, nil)
}
