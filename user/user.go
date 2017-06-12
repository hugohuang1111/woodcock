package user

import (
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/hugohuang1111/poker/db"
	"github.com/hugohuang1111/woodcock/constants"
)

var (
	guestCreateMutex sync.Mutex
	activityUsers    map[uint64]uint64
)

func init() {
	activityUsers = make(map[uint64]uint64)
}

func getNano() string {
	var nanos int64
	guestCreateMutex.Lock()
	now := time.Now()
	nanos = now.UnixNano()
	guestCreateMutex.Unlock()

	return string(nanos)
}

func register(connectID uint64, msg []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1

	name, _ := jsonparser.GetUnsafeString(msg, "name")
	passwd, _ := jsonparser.GetUnsafeString(msg, "password")
	if 0 == len(name) {
		name = "user" + getNano()
	}
	if 0 == len(passwd) {
		passwd = "111111"
	}

	suc, _ := db.UserRegister(name, passwd)

	if suc {
		activityUsers[connectID] = 0 //TODO uid
		constants.SetRespError(resp, constants.ERROR_SUCCESS)
		resp["user"] = name
		resp["passwd"] = passwd
	} else {
		constants.SetRespError(resp, constants.ERROR_REGISTER_FAIL)
	}

	return resp
}

func login(connectID uint64, msg []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1

	name, eName := jsonparser.GetUnsafeString(msg, "name")
	pw, ePW := jsonparser.GetUnsafeString(msg, "password")
	if nil != eName || nil != ePW {
		constants.SetRespError(resp, constants.ERROR_PARAM_NIL)
		return resp
	}

	if suc, uid := db.UserLogin(name, pw); suc {
		constants.SetRespError(resp, constants.ERROR_SUCCESS)
		activityUsers[connectID] = uid
	} else {
		constants.SetRespError(resp, constants.ERROR_LOGIN_FAIL)
	}

	return resp
}

func logout(connectID uint64, msg []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	delete(activityUsers, connectID)

	return resp
}
