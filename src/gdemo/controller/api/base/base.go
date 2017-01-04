package base

import (
	"andals/gobox/http/controller"
	glog "andals/gobox/log"

	"gdemo/log"
	"gdemo/misc"
)

func RegAction(cl *controller.Controller) {
	cl.AddBeforeAction(".+", BeforeAction)
	cl.AddAfterAction(".+", AfterAction)
	cl.AddDestructFunc(".+", Destruct)
}

func BeforeAction(context *controller.Context, args []string) {
	realIp := context.RemoteRealAddr.Ip

	var logger glog.ILogger
	sl, e := glog.NewSimpleLogger(log.UserLogWriter, glog.LEVEL_INFO, glog.NewWebFormater(context.Rid, []byte(realIp)))
	if e != nil {
		logger = new(glog.NoopLogger)
	} else {
		logger = glog.NewAsyncLogger(sl, log.UserLogRoutineCh)
	}

	context.TransData[misc.TRANS_DATA_KEY_USER_LOG] = logger

	context.RespWriter.Header().Add("X-Powered-By", "gohttp")
}

func AfterAction(context *controller.Context, args []string) {

}

func Destruct(transData map[string]interface{}, args []string) {
	l, ok := transData[misc.TRANS_DATA_KEY_USER_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Free()
		}
	}
}
