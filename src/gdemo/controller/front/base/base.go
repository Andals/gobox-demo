package base

import (
	"andals/gobox/http/controller"
	glog "andals/gobox/log"

	"family/log"
	sathena "family/svc/athena"
)

func RegAction(cl *controller.Controller) {
	cl.AddBeforeAction(".+", BeforeAction)
	cl.AddAfterAction(".+", AfterAction)
	cl.AddDestructFunc(".+", destruct)
}

func BeforeAction(context *controller.Context, args []string) {
	realIp := context.RemoteRealAddr.Ip
	if IsHttps(context) {
		realIp = context.Req.Header.Get("HTTPS-REAL-IP")
	}
	context.TransData[sathena.TRANS_DATA_KEY_REAL_IP] = realIp

	var logger glog.ILogger
	sl, e := glog.NewSimpleLogger(log.UserLogWriter, glog.LEVEL_INFO, glog.NewWebFormater(context.Rid, []byte(realIp)))
	if e != nil {
		logger = new(glog.NoopLogger)
	} else {
		logger = glog.NewAsyncLogger(sl, log.UserLogRoutineCh)
	}

	context.TransData[sathena.TRANS_DATA_KEY_USER_LOG] = logger

	logger, e = glog.NewSimpleLogger(log.FrontAthenaLogWriter, glog.LEVEL_INFO, glog.NewWebFormater(context.Rid, []byte(realIp)))
	if e != nil {
		logger = new(glog.NoopLogger)
	}
	context.TransData[sathena.TRANS_DATA_KEY_FRONT_ATHENA_LOG] = logger

	context.RespWriter.Header().Add("X-Powered-By", "golang")
}

func AfterAction(context *controller.Context, args []string) {
	l, ok := context.TransData[sathena.TRANS_DATA_KEY_USER_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Free()
		}
	}

	l, ok = context.TransData[sathena.TRANS_DATA_KEY_FRONT_ATHENA_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Free()
		}
	}
}

func destruct(transData map[string]interface{}, args []string) {
	l, ok := transData[sathena.TRANS_DATA_KEY_USER_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Free()
		}
	}

	l, ok = transData[sathena.TRANS_DATA_KEY_FRONT_ATHENA_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Free()
		}
	}
}

func IsHttps(context *controller.Context) bool {
	if context.Req.Header.Get("HTTPS-FLAG") != "" {
		return true
	}
	return false
}
