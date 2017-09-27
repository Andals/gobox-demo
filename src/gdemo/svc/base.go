package svc

import (
	"andals/golog"
)

type baseSvc struct {
	elogger golog.ILogger
}

func newBaseSvc(elogger golog.ILogger) *baseSvc {
	if elogger == nil {
		elogger = new(golog.NoopLogger)
	}

	return &baseSvc{
		elogger: elogger,
	}
}
