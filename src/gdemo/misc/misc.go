package misc

import (
	"encoding/json"

	"andals/gobox/exception"

	"gdemo/errno"
)

type ApiData struct {
	Errno int `json:"errno"`

	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ApiJson(data interface{}, e *exception.Exception) []byte {
	result := &ApiData{
		Errno: errno.SUCCESS,
		Msg:   "",
		Data:  data,
	}
	if e != nil {
		result.Errno = e.Errno()
		result.Msg = e.Msg()
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
