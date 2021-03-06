package demo

import (
	"gdemo/errno"
	"gdemo/svc"

	"github.com/andals/gobox/exception"
	"github.com/andals/gobox/http/query"
)

func (this *DemoController) AddAction(context *DemoContext) {
	ap, e := this.parseAddActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	ids, err := context.demoSvc.Insert(ap)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_API_DEMO_INSERT_FAILED, err.Error())
		return
	}

	context.ApiData.Data = map[string]interface{}{
		"id": ids[0],
	}
}

func (this *DemoController) parseAddActionParams(context *DemoContext) (*svc.DemoEntity, *exception.Exception) {
	ap := new(svc.DemoEntity)

	qs := query.NewQuerySet()
	qs.StringVar(&ap.Name, "name", true, errno.E_API_DEMO_INVALID_NAME, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", true, errno.E_API_DEMO_INVALID_STATUS, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, e
	}

	if ap.Status < 0 {
		return ap, exception.New(errno.E_API_DEMO_INVALID_STATUS, "invalid status")
	}

	return ap, nil
}
