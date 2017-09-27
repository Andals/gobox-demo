package demo

import (
	"gdemo/controller/api"
	"gdemo/svc"

	gcontroller "andals/gobox/http/controller"
	"gdemo/conf"
)

type DemoController struct {
	api.BaseController

	demoSvc *svc.DemoSvc
}

func (this *DemoController) BeforeAction(context gcontroller.ActionContext) {
	this.BaseController.BeforeAction(context)

	acontext := context.(*api.ApiContext)
	this.demoSvc = svc.NewDemoSvc(
		acontext.ErrorLogger,
		acontext.MysqlClient,
		conf.BaseConf.PrjName,
		acontext.RedisClient,
	)
}
