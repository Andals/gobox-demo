package misc

import (
	"encoding/json"

	"andals/gobox/exception"

	"gdemo/errno"
)

const (
	TRANS_DATA_KEY_USER_LOG = "user_log"
)

type ApiData struct {
	Errno int `json:"errno"`

	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ApiJson(data interface{}, e *exception.Exception) []byte {
	en := errno.SUCCESS
	msg := ""
	if e != nil {
		en = e.Errno()
		msg = e.Msg()
	}

	result := &ApiData{
		Errno: en,
		Msg:   msg,
		Data:  data,
	}

	aj, err := json.Marshal(result)
	if err != nil {
		result.Errno = errno.E_COMMON_JSON_ENCODE_ERROR
		result.Msg = err.Error()
		result.Data = nil

		aj, _ = json.Marshal(result)
	}

	return aj
}
