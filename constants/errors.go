package constants

const (
	ERROR_SUCCESS = 0
	//client error
	ERROR_MSG_FORMAT_ERROR = 1000
	ERROR_REGISTER_FAIL    = 1001
	ERROR_PARAM_NIL        = 1002
	ERROR_LOGIN_FAIL       = 1003
	ERROR_PARAM_WRONG      = 1004
	ERROR_NOT_AUTHOR       = 1005
)

var (
	ErrorMap = map[int]string{
		ERROR_MSG_FORMAT_ERROR: "msg format error",
		ERROR_REGISTER_FAIL:    "register failed",
		ERROR_PARAM_NIL:        "param is nil",
		ERROR_LOGIN_FAIL:       "login failed",
		ERROR_PARAM_WRONG:      "param wrong",
		ERROR_NOT_AUTHOR:       "not authorize",
	}
)

func GenErrorMsg(e int) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	resp["error"] = e
	resp["description"] = ErrorMap[e]

	return resp
}

func SetRespError(resp map[string]interface{}, e int) map[string]interface{} {
	resp["error"] = e
	resp["description"] = ErrorMap[e]

	return resp
}
