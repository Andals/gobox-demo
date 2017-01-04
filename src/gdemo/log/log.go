package log

import (
	"gdemo/errno"

	"andals/gobox/exception"
	glog "andals/gobox/log"
	"andals/gobox/log/buffer"
	"andals/gobox/log/writer"
)

const (
	USER_LOG_ROUTINE_QUEUE_LEN = 4096
	USER_LOG_BUFFER_SIZE       = 4096
)

var UserLogWriter *buffer.Buffer
var UserLogRoutineCh *glog.AsyncLogRoutineCh

func Init(logRoot string) *exception.Exception {
	fw, err := writer.NewFileWriter(logRoot + "/user.log")
	if err != nil {
		return exception.New(errno.E_SYS_INIT_LOG_FAIL, err.Error())
	}

	UserLogWriter = buffer.NewBuffer(fw, USER_LOG_BUFFER_SIZE)
	UserLogRoutineCh = glog.NewAsyncLogRoutine(USER_LOG_ROUTINE_QUEUE_LEN)

	return nil
}

func Free() {
	glog.FreeAsyncLogRoutines()
}
