package api

import (
	"andals/gobox/http/controller"

	"gdemo/controller/api/base"
	"gdemo/controller/api/index"
)

func RegAction(cl *controller.Controller) {
	base.RegAction(cl)
	index.RegAction(cl)
}
