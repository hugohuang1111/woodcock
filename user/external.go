package user

import (
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/constants"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

//Module gate module
type Module struct {
}

//NewModule new gate module
func NewModule() *Module {
	m := new(Module)

	return m
}

//OnInit module init
func (m *Module) OnInit() {
	glog.Info("user init")
}

//OnDestroy module destroy
func (m *Module) OnDestroy() {
}

//OnEvent module event
func (m *Module) OnEvent(msg *module.Message) {
	connID := module.GetConnectID(msg.Payload)
	clientData := module.GetClientData(msg.Payload)
	switch msg.Type {
	case module.MsgTypeDisconnect:
		m := new(module.Message)
		m.Recver = module.MOD_ROOM
		m.Sender = module.MOD_USER
		m.Type = module.MsgTypeDisconnect
		if nil == m.Payload {
			m.Payload = make(map[string]interface{})
		}
		m.Payload[module.PayloadKeyUserID] = activityUsers[connID]
		router.Route(m)

		logout(connID, clientData)
	case module.MsgTypeGetUserID:
		{
			uid := getUIDByConnID(connID)
			msg.Payload[module.PayloadKeyUserID] = uid
			m.response(msg, nil)
		}
	case module.MsgTypeClient:
		{
			v, e := jsonparser.GetString(clientData, "type")
			if nil != e {
				m.response(msg, constants.GenErrorMsg(constants.ERROR_MSG_FORMAT_ERROR))
				return
			}

			s := strings.Split(v, ":")
			cmd := s[1]
			var resp map[string]interface{}
			switch cmd {
			case "register":
				{
					resp = register(connID, clientData)
				}
			case "login":
				{
					resp = login(connID, clientData)
				}
			case "logout":
				{
					resp = logout(connID, clientData)
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
			glog.Error("user unknow msg type:", msg.Type)
		}
	}
}

func (m *Module) response(req *module.Message, payload map[string]interface{}) {
	msg := new(module.Message)
	msg.Recver = req.Sender
	msg.Sender = module.MOD_USER
	msg.Type = req.Type
	msg.Payload = req.Payload

	if nil != payload {
		jsonString, err := json.Marshal(payload)
		if nil != err {
			glog.Error("marshal resp to string failed:", err)
			return
		}
		msg.Payload[module.PayloadKeyClientData] = jsonString
	}

	router.Route(msg)
}

func (m *Module) handler(cmd string) {

}
