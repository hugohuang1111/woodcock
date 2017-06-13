package room

import (
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/constants"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

//Module room module
type Module struct {
}

//NewModule new room module
func NewModule() *Module {
	m := new(Module)

	return m
}

//OnInit module init
func (m *Module) OnInit() {
	glog.Info("room init")
}

//OnDestroy module destroy
func (m *Module) OnDestroy() {
}

//OnEvent module event
func (m *Module) OnEvent(msg *module.Message) {
	switch msg.Type {
	case module.MOD_MSG_TYPE_DISCONNECT:
		{
			leave(msg.ConnectID, msg.Payload)
		}
	case module.MOD_MSG_TYPE_GET_USER_ID:
		if 0 == msg.Userid {
			m.response(msg, constants.GenErrorMsg(constants.ERROR_NOT_AUTHOR))
			return
		}
		fallthrough
	case module.MOD_MSG_TYPE_CLIENT:
		{
			t, e := jsonparser.GetString(msg.Payload, "type")
			if nil != e {
				m.response(msg, constants.GenErrorMsg(constants.ERROR_MSG_FORMAT_ERROR))
				return
			}

			if 0 == msg.Userid {
				m.forward(msg, module.MOD_USER)
				return
			}

			s := strings.Split(t, ":")
			cmd := s[1]
			var resp map[string]interface{}
			switch cmd {
			case "entry":
				{
					resp = entry(msg.Userid, msg.Payload)
				}
			case "leave":
				{
					resp = leave(msg.Userid, msg.Payload)
				}
			default:
				{
					glog.Error("user unknow cmd:", cmd)
				}
			}

			if nil != resp {
				m.response(msg, resp)
			}
		}
	default:
		{
			glog.Error("room unknow msg type:", msg.Type)
		}
	}
}

func (m *Module) response(req *module.Message, payload map[string]interface{}) {
	msg := new(module.Message)
	msg.ConnectID = req.ConnectID
	msg.Sender = module.MOD_USER
	msg.Recver = req.Sender
	msg.Userid = req.Userid
	msg.Type = req.Type

	jsonString, err := json.Marshal(payload)
	if nil != err {
		glog.Error("marshal resp to string failed:", err)
		return
	}
	msg.Payload = jsonString

	router.Route(msg)
}

func (m *Module) forward(msg *module.Message, receiver string) {
	msg.Sender = module.MOD_ROOM
	msg.Recver = receiver
	msg.Type = module.MOD_MSG_TYPE_GET_USER_ID
	router.Route(msg)
}
