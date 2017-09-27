package demo

import (
	"gdemo/controller/api"
	"gdemo/errno"

	"andals/gobox/exception"
	"andals/gobox/http/query"
)

type getActionParams struct {
	id int64
}

func (this *DemoController) GetAction(context *api.ApiContext) {
	ap, e := this.parseGetActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	entity, err := this.demoSvc.GetById(ap.id)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		return
	}

	context.ApiData.Data = entity
}

func (this *DemoController) parseGetActionParams(context *api.ApiContext) (*getActionParams, *exception.Exception) {
	ap := new(getActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.E_COMMON_INVALID_ID, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
