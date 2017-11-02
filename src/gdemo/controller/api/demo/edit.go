package demo

import (
	"gdemo/errno"
	"gdemo/svc"

	"andals/gobox/exception"
	"andals/gobox/http/query"
)

func (this *DemoController) EditAction(context *DemoContext) {
	ap, exists, e := this.parseEditActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	if len(exists) == 1 { //only has id
		return
	}

	updated, err := context.demoSvc.UpdateById(ap.Id, ap, exists)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_API_DEMO_UPDATE_FAILED, err.Error())
		return
	}

	context.ApiData.Data = updated
}

func (this *DemoController) parseEditActionParams(context *DemoContext) (*svc.DemoEntity, map[string]bool, *exception.Exception) {
	ap := new(svc.DemoEntity)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.Id, "id", true, errno.E_COMMON_INVALID_ID, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.Name, "name", false, errno.E_API_DEMO_INVALID_NAME, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", false, errno.E_API_DEMO_INVALID_STATUS, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, exception.New(errno.E_API_DEMO_INVALID_STATUS, "invalid status")
	}

	return ap, qs.ExistsInfo(), nil
}
