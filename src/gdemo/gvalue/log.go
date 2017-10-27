package gvalue

import (
	"gdemo/conf"
	"gdemo/errno"

	"andals/gobox/exception"
	"andals/golog"
)

var ErrorLogWriter golog.IWriter
var RedisLogWriter golog.IWriter
var MysqlLogWriter golog.IWriter

func InitLog(systemName string) *exception.Exception {
	golog.InitBufferAutoFlushRoutine(conf.LogConf.MaxBufferNum, conf.LogConf.BufferAutoFlushTimeInterval)
	golog.InitAsyncLogRoutine(conf.LogConf.MaxAsyncMsgNum)

	var err error

	ErrorLogWriter, err = golog.NewFileWriter(conf.LogConf.RootPath + "/" + systemName + "_error.log")
	if err != nil {
		return exception.New(errno.E_SYS_INIT_LOG_FAIL, err.Error())
	}

	RedisLogWriter, err = golog.NewFileWriter(conf.LogConf.RootPath + "/" + systemName + "_redis.log")
	if err != nil {
		return exception.New(errno.E_SYS_INIT_LOG_FAIL, err.Error())
	}

	MysqlLogWriter, err = golog.NewFileWriter(conf.LogConf.RootPath + "/" + systemName + "_mysql.log")
	if err != nil {
		return exception.New(errno.E_SYS_INIT_LOG_FAIL, err.Error())
	}

	return nil
}

func NewSyncLogger(w golog.IWriter, formater golog.IFormater) golog.ILogger {
	bw := golog.NewBuffer(w, conf.LogConf.Bufsize)
	logger, err := golog.NewSimpleLogger(bw, conf.LogConf.Level, formater)
	if err != nil {
		return new(golog.NoopLogger)
	}

	return logger
}

func NewAsyncLogger(w golog.IWriter, formater golog.IFormater) golog.ILogger {
	var logger golog.ILogger
	bw := golog.NewBuffer(w, conf.LogConf.Bufsize)
	sl, e := golog.NewSimpleLogger(bw, conf.LogConf.Level, formater)
	if e != nil {
		logger = new(golog.NoopLogger)
	} else {
		logger = golog.NewAsyncLogger(sl)
	}

	return logger
}

func FreeLog() {
	golog.FreeAsyncLogRoutine()
	golog.FreeBuffers()
}
