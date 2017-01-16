package api

import (
	"net/http"

	"andals/gobox/exception"
	gcontroller "andals/gobox/http/controller"
	glog "andals/gobox/log"

	"gdemo/controller"
	"gdemo/log"
	"gdemo/misc"
)

type ApiContext struct {
	*controller.BaseContext

	ApiData struct {
		Data interface{}
		Err  *exception.Exception
	}

	UserLogger glog.ILogger
}

type BaseController struct {
	controller.BaseController
}

func (this *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(ApiContext)
	context.BaseContext = this.BaseController.NewActionContext(req, respWriter).(*controller.BaseContext)

	return context
}

func (this *BaseController) BeforeAction(context gcontroller.ActionContext) {
	acontext := context.(*ApiContext)
	realIp := acontext.RemoteRealAddr.Ip

	sl, e := glog.NewSimpleLogger(log.UserLogWriter, glog.LEVEL_INFO, glog.NewWebFormater(acontext.Rid, []byte(realIp)))
	if e != nil {
		acontext.UserLogger = new(glog.NoopLogger)
	} else {
		acontext.UserLogger = glog.NewAsyncLogger(sl, log.UserLogRoutineCh)
	}

	acontext.RespWriter.Header().Add("X-Powered-By", "gohttp")
}

func (this *BaseController) AfterAction(context gcontroller.ActionContext) {
	acontext := context.(*ApiContext)

	acontext.RespBody = misc.ApiJson(acontext.ApiData.Data, acontext.ApiData.Err)
}

func (this *BaseController) Destruct(context gcontroller.ActionContext) {
	acontext := context.(*ApiContext)

	acontext.UserLogger.Free()
}
