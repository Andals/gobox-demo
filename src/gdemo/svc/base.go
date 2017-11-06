package svc

import (
	"github.com/andals/golog"
)

type BaseSvc struct {
	elogger golog.ILogger
}

func NewBaseSvc(elogger golog.ILogger) *BaseSvc {
	if elogger == nil {
		elogger = new(golog.NoopLogger)
	}

	return &BaseSvc{
		elogger: elogger,
	}
}
