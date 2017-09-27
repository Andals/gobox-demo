package controller

import (
	"gdemo/gvalue"

	"andals/gobox/encoding"
	"andals/gobox/http/controller"
	gmisc "andals/gobox/misc"
	"andals/golog"

	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"

	DOWNSTREAM_SERVER_IP = "127.0.0.1"
)

type BaseContext struct {
	Req        *http.Request
	RespWriter http.ResponseWriter
	RespBody   []byte

	QueryValues    url.Values
	RemoteRealAddr struct {
		Ip   string
		Port string
	}
	Rid []byte

	LogFormater golog.IFormater
	ErrorLogger golog.ILogger
}

func (this *BaseContext) Request() *http.Request {
	return this.Req
}

func (this *BaseContext) ResponseWriter() http.ResponseWriter {
	return this.RespWriter
}

func (this *BaseContext) ResponseBody() []byte {
	return this.RespBody
}

type BaseController struct {
}

func (this *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	context := &BaseContext{
		Req:        req,
		RespWriter: respWriter,
	}

	req.ParseForm()
	context.QueryValues = req.Form

	context.RemoteRealAddr.Ip, context.RemoteRealAddr.Port = this.parseRemoteAddr(req)

	now := time.Now()
	timeInt := now.UnixNano()
	randInt := gmisc.RandByTime(&now)

	ridStr := context.RemoteRealAddr.Ip + ":" + context.RemoteRealAddr.Port + "," + strconv.FormatInt(timeInt, 10) + "," + strconv.FormatInt(randInt, 10)
	context.Rid = encoding.Base64Encode([]byte(ridStr))

	context.LogFormater = golog.NewWebFormater(context.Rid, []byte(context.RemoteRealAddr.Ip))
	context.RespWriter.Header().Add("X-Powered-By", "gohttp")
	context.ErrorLogger = gvalue.NewAsyncLogger(gvalue.ErrorLogWriter, context.LogFormater)

	return context
}

func (this *BaseController) BeforeAction(context controller.ActionContext) {
}

func (this *BaseController) AfterAction(context controller.ActionContext) {
}

func (this *BaseController) Destruct(context controller.ActionContext) {
	bcontext := context.(*BaseContext)

	bcontext.RespBody = nil
	bcontext.QueryValues = nil
	bcontext.Rid = nil

	bcontext.ErrorLogger.Free()
	bcontext.LogFormater = nil
}

func (this *BaseController) parseRemoteAddr(req *http.Request) (string, string) {
	rs := strings.Split(req.RemoteAddr, ":")
	if rs[0] == DOWNSTREAM_SERVER_IP {
		ip := strings.TrimSpace(req.Header.Get(REMOTE_REAL_IP_HEADER_KEY))
		port := strings.TrimSpace(req.Header.Get(REMOTE_REAL_PORT_HEADER_KEY))
		if ip != "" && port != "" {
			return ip, port
		}
	}

	return rs[0], rs[1]
}
