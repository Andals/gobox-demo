package front

import (
	"andals/gobox/http/controller"

	"family/controller/front/athena"
	"family/controller/front/base"
)

func RegAction(cl *controller.Controller) {
	base.RegAction(cl)
	athena.RegAction(cl)
}
