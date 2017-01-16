package api

import (
	"andals/gobox/http/query"

	"gdemo/errno"
)

type IndexController struct {
	BaseController
}

func (this *IndexController) IndexAction(context *ApiContext) {
	var a int
	var b string

	qs := query.NewQuerySet()
	qs.IntVar(&a, "a", errno.E_API_INVALID_A, "invalid a", nil)
	qs.StringVar(&b, "b", errno.E_API_INVALID_B, "invalid b", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	context.UserLogger.Info(context.Rid)

	var data struct {
		Rid []byte
	}
	data.Rid = context.Rid
	context.ApiData.Data = data
}
