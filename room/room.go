package room

import (
	"github.com/buger/jsonparser"
	"github.com/hugohuang1111/woodcock/constants"
)

func entry(uid uint64, payload []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	t, _ := jsonparser.GetUnsafeString(payload, "type")
	resp["type"] = t

	roomID, _ := jsonparser.GetInt(payload, "roomID")
	if 0 == roomID {
		constants.SetRespError(resp, constants.ERROR_PARAM_WRONG)
	} else {

	}

	return resp
}

func leave(connID uint64, payload []byte) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["version"] = 1
	t, eType := jsonparser.GetUnsafeString(payload, "type")
	if nil == eType {
		resp["type"] = t
	}

	//TODO get user id by connID

	return resp
}
