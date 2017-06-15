package user

import (
	"strconv"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/hugohuang1111/woodcock/constants"
	"github.com/hugohuang1111/woodcock/db"
)

var (
	guestCreateMutex sync.Mutex
	activityUsers    map[uint64]uint64
)

func init() {
	activityUsers = make(map[uint64]uint64) //connID->userID
}

func getNano() string {
	var nanos int64
	guestCreateMutex.Lock()
	now := time.Now()
	nanos = now.UnixNano()
	guestCreateMutex.Unlock()

	return strconv.FormatInt(nanos, 10)
}

func register(connectID uint64, msg []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	t, eType := jsonparser.GetUnsafeString(msg, "type")
	if nil == eType {
		resp["type"] = t
	}

	name, _ := jsonparser.GetUnsafeString(msg, "name")
	passwd, _ := jsonparser.GetUnsafeString(msg, "password")
	if 0 == len(name) {
		name = "user" + getNano()
	}
	if 0 == len(passwd) {
		passwd = "111111"
	}

	uid, err := db.UserRegister(name, passwd)

	if nil == err {
		activityUsers[connectID] = uid
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
	t, eType := jsonparser.GetUnsafeString(msg, "type")
	if nil == eType {
		resp["type"] = t
	}

	name, eName := jsonparser.GetUnsafeString(msg, "name")
	pw, ePW := jsonparser.GetUnsafeString(msg, "password")
	if nil != eName || nil != ePW {
		constants.SetRespError(resp, constants.ERROR_PARAM_NIL)
		return resp
	}

	if suc, uid := db.UserLogin(name, pw); suc {
		constants.SetRespError(resp, constants.ERROR_SUCCESS)
		activityUsers[connectID] = uid
		resp["userID"] = uid
	} else {
		constants.SetRespError(resp, constants.ERROR_LOGIN_FAIL)
	}

	return resp
}

func logout(connectID uint64, msg []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	t, eType := jsonparser.GetUnsafeString(msg, "type")
	if nil == eType {
		resp["type"] = t
	}

	delete(activityUsers, connectID)

	return resp
}

func getUIDByConnID(connID uint64) uint64 {
	uid, _ := activityUsers[connID]

	return uid
}
