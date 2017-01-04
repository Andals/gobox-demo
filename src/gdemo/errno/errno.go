package errno

const (
	SUCCESS = 0

	E_SYS_INVALID_FLAG_PRJ_HOME  = 11
	E_SYS_SAVE_PID_FILE_FAIL     = 12
	E_SYS_INIT_LOG_FAIL          = 13

	E_COMMON_JSON_ENCODE_ERROR = 101
	E_COMMON_FILE_NOT_EXIST    = 102
	E_COMMON_READ_FILE_ERROR   = 103
	E_COMMON_JSON_DECODE_ERROR = 104

	E_CONF_INVALID_PRJ_HOME    = 201
	E_CONF_INVALID_SERVER_CONF = 202
)