package demo

import (
	"gdemo/controller/api"
	"gdemo/errno"

	"andals/gobox/exception"
	"andals/gobox/http/query"
)

type delActionParams struct {
	id int64
}

func (this *DemoController) DelAction(context *api.ApiContext) {
	ap, e := this.parseDelActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	deleted, err := this.demoSvc.DeleteById(ap.id)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		return
	}

	context.ApiData.Data = deleted
}

func (this *DemoController) parseDelActionParams(context *api.ApiContext) (*delActionParams, *exception.Exception) {
	ap := new(delActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.E_COMMON_INVALID_ID, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
