package room

import (
	"github.com/buger/jsonparser"
	"github.com/hugohuang1111/woodcock/constants"
)

type userInfo struct {
	connID uint64
	roomID uint64
}

var (
	userConnMap map[uint64]userInfo //userID->connectID
)

func init() {
	userConnMap = make(map[uint64]userInfo)
}

func entry(connID, uid, rid uint64, payload []byte) map[string]interface{} {
	var roomID uint64
	var clientDataType string
	if nil != payload {
		intID, _ := jsonparser.GetInt(payload, "roomID")
		roomID = uint64(intID)
		clientDataType, _ = jsonparser.GetUnsafeString(payload, "type")
	} else {
		roomID = rid
	}

	var bSitDown bool
	if 0 != roomID {
		t := getOrCreateTable(roomID)
		bSitDown = t.sitDown(uid)
	}

	resp := make(map[string]interface{})
	resp["version"] = 1
	resp["type"] = clientDataType

	if !bSitDown {
		constants.SetRespError(resp, constants.ERROR_ROOM_FULL)
	} else if 0 == roomID {
		constants.SetRespError(resp, constants.ERROR_PARAM_WRONG)
	} else {
		constants.SetRespError(resp, constants.ERROR_SUCCESS)
		userConnMap[uid] = userInfo{connID, roomID}
	}

	return resp
}

func leave(connID, uID uint64, payload []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	t, eType := jsonparser.GetUnsafeString(payload, "type")
	if nil == eType {
		resp["type"] = t
	}

	delete(userConnMap, uID)
	return resp
}

func abandonSuit(uID uint64, payload []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	t, eType := jsonparser.GetUnsafeString(payload, "type")
	if nil == eType {
		resp["type"] = t
	}
	suit, eType := jsonparser.GetInt(payload, "suit")

	var e error
	if info, ok := userConnMap[uID]; ok {
		r := getOrCreateTable(info.roomID)
		e = r.abandonSuit(uID, int(suit))
	}
	if nil == e {
		constants.SetRespError(resp, constants.ERROR_SUCCESS)
	} else {
		constants.SetRespError(resp, constants.ERROR_ABANDON_SUIT_FAIL)
	}

	return resp

}
