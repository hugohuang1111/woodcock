package module

import (
	"strings"

	"github.com/buger/jsonparser"
)

//Message message pass by modules
type Message struct {
	Sender  string
	Recver  string
	Payload map[string]interface{}
	Type    string
}

func GetConnectID(payload map[string]interface{}) uint64 {
	return GetUint64(payload, PayloadKeyConnectID)
}

func GetUserID(payload map[string]interface{}) uint64 {
	return GetUint64(payload, PayloadKeyUserID)
}

func GetClientData(payload map[string]interface{}) []byte {
	if v, exist := payload[PayloadKeyClientData]; exist {
		if data, ok := v.([]byte); ok {
			return data
		}
	}
	return nil
}

func GetClientDataCmd(clientData []byte) string {
	if t, e := jsonparser.GetString(clientData, "type"); nil == e {
		if s := strings.Split(t, ":"); 2 == len(s) {
			return s[1]
		}
	}
	return ""
}

func GetUint64(payload map[string]interface{}, key string) uint64 {
	v, exist := payload[key]
	if !exist {
		return 0
	}
	if id, ok := v.(uint64); ok {
		return id
	}
	return 0
}
