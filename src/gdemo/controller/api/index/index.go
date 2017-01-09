package index

import (
	"andals/gobox/http/controller"
	glog "andals/gobox/log"

	"andals/gobox/http/query"
	"gdemo/controller/api"
	"gdemo/errno"
	"gdemo/misc"
)

func RegAction(cl *controller.Controller) {
	api.RegAction(cl)

	cl.AddExactMatchAction("/api/index", Index)

	api.Foo()
}

func Index(context *controller.Context, args []string) {
	var a int
	var b string

	qs := query.NewQuerySet()
	qs.IntVar(&a, "a", errno.E_API_INVALID_A, "invalid a", nil)
	qs.StringVar(&b, "b", errno.E_API_INVALID_B, "invalid b", nil)
	//e := qs.Parse(context.QueryValues)

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
