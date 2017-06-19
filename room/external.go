package room

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/constants"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

//Module room module
type Module struct {
	skelection *module.Skelecton
}

//NewModule new room module
func NewModule() *Module {
	m := new(Module)

	return m
}

//OnInit module init
func (m *Module) OnInit() {
	glog.Info("room init")
	m.skelection.Run(m)
}

//OnDestroy module destroy
func (m *Module) OnDestroy() {
	m.skelection.Stop()
}

//OnEvent module event
func (m *Module) OnEvent(msg *module.Message) {
	m.skelection.Add(msg)
}

//OnMsg module skelection interface
func (m *Module) OnMsg(msg *module.Message) {
	connID := module.GetConnectID(msg.Payload)
	uID := module.GetUserID(msg.Payload)
	clientData := module.GetClientData(msg.Payload)
	switch msg.Type {
	case module.MsgTypeDisconnect:
		leave(connID, uID, clientData)
	case module.MsgTypeEntryRoom:
		rID := module.GetUint64(msg.Payload, module.PayloadKeyRoomID)
		entry(connID, uID, rID, nil)
	case module.MsgTypeGetUserID:
		if 0 == uID {
			message := new(module.Message)
			message.Recver = module.MOD_GATE
			message.Sender = module.MOD_ROOM
			message.Type = module.MsgTypeClient
			message.Payload = make(map[string]interface{})
			message.Payload[module.PayloadKeyUserID] = uID
			message.Payload[module.PayloadKeyConnectID] = connID
			resp := constants.GenErrorMsg(constants.ERROR_NOT_AUTHOR)
			t, e := jsonparser.GetString(clientData, "type")
			if nil == e {
				resp["type"] = t
			}
			m.send(message, resp)
			return
		}
		msg.Sender = module.MOD_GATE
		msg.Type = module.MsgTypeClient
		fallthrough
	case module.MsgTypeClient:
		cmd := module.GetClientDataCmd(clientData)
		if 0 == len(cmd) {
			m.response(msg, constants.GenErrorMsg(constants.ERROR_MSG_FORMAT_ERROR))
			return
		}
		if 0 == uID {
			m.forward(msg, module.MOD_USER)
			return
		}

		var resp map[string]interface{}
		switch cmd {
		case "entry":
			resp = entry(connID, uID, 0, clientData)
		case "leave":
			resp = leave(connID, uID, clientData)
		case "abanbonSuit":
			uID := module.GetUserID(msg.Payload)
			resp = leave(connID, uID, clientData)
		default:
			glog.Error("user unknow cmd:", cmd)
		}

		if nil != resp {
			m.response(msg, resp)
		}
	default:
		glog.Error("room unknow msg type:", msg.Type)
	}

}

func (m *Module) send(msg *module.Message, payload map[string]interface{}) {
	jsonString, err := json.Marshal(payload)
	if nil != err {
		glog.Error("marshal resp to string failed:", err)
		return
	}
	msg.Payload[module.PayloadKeyClientData] = jsonString

	router.Route(msg)
}

func (m *Module) response(req *module.Message, payload map[string]interface{}) {
	msg := new(module.Message)
	msg.Sender = module.MOD_ROOM
	msg.Recver = req.Sender
	msg.Type = req.Type
	msg.Payload = req.Payload

	jsonString, err := json.Marshal(payload)
	if nil != err {
		glog.Error("marshal resp to string failed:", err)
		return
	}
	msg.Payload[module.PayloadKeyClientData] = jsonString

	router.Route(msg)
}

func (m *Module) forward(msg *module.Message, receiver string) {
	msg.Sender = module.MOD_ROOM
	msg.Recver = receiver
	msg.Type = module.MsgTypeGetUserID
	router.Route(msg)
}
