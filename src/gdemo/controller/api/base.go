package api

import (
	"andals/gobox/http/system"
	"andals/gobox/mysql"
	"gdemo/controller"
	"gdemo/errno"
	"gdemo/gvalue"
	"gdemo/misc"

	"andals/gobox/exception"
	gcontroller "andals/gobox/http/controller"
	"andals/gobox/redis"
	"andals/golog"

	"html"
	"net/http"
)

type ApiContext struct {
	*controller.BaseContext

	ApiData struct {
		V    string
		Data interface{}
		Err  *exception.Exception
	}

	MysqlPool   *mysql.Pool
	MysqlClient *mysql.Client
	MysqlLogger golog.ILogger

	RedisPool   *redis.Pool
	RedisClient *redis.Client
	RedisLogger golog.ILogger
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
	this.BaseController.BeforeAction(acontext.BaseContext)

	var err error
	acontext.MysqlPool = gvalue.MysqlClientPool
	acontext.MysqlClient, err = acontext.MysqlPool.Get()
	if err != nil {
		acontext.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		system.JumpOutAction(this.jumpToError)
	}
	acontext.MysqlLogger = gvalue.NewAsyncLogger(gvalue.MysqlLogWriter, acontext.LogFormater)
	acontext.MysqlClient.SetLogger(acontext.MysqlLogger)

	acontext.RedisPool = gvalue.RedisClientPool
	acontext.RedisClient, err = acontext.RedisPool.Get()
	if err != nil {
		acontext.ApiData.Err = exception.New(errno.E_SYS_REDIS_ERROR, err.Error())
		system.JumpOutAction(this.jumpToError)
	}
	acontext.RedisLogger = gvalue.NewAsyncLogger(gvalue.RedisLogWriter, acontext.LogFormater)
	acontext.RedisClient.SetLogger(acontext.RedisLogger)
}

func (this *BaseController) AfterAction(context gcontroller.ActionContext) {
	acontext := context.(*ApiContext)

	f := acontext.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := acontext.QueryValues.Get("_callback")
		if callback != "" {
			acontext.RespBody = misc.ApiJsonp(acontext.ApiData.V, acontext.ApiData.Data, acontext.ApiData.Err, html.EscapeString(callback))
			return
		}
	}

	acontext.RespBody = misc.ApiJson(acontext.ApiData.V, acontext.ApiData.Data, acontext.ApiData.Err)

	this.BaseController.AfterAction(acontext.BaseContext)
}

func (this *BaseController) Destruct(context gcontroller.ActionContext) {
	acontext := context.(*ApiContext)

	acontext.MysqlPool.Put(acontext.MysqlClient)
	acontext.MysqlLogger.Free()

	acontext.RedisPool.Put(acontext.RedisClient)
	acontext.RedisLogger.Free()

	this.BaseController.Destruct(acontext.BaseContext)
}

func (this *BaseController) jumpToError(context gcontroller.ActionContext, args ...interface{}) {
	acontext := context.(*ApiContext)

	acontext.RespBody = misc.ApiJson(acontext.ApiData.V, acontext.ApiData.Data, acontext.ApiData.Err)
}
