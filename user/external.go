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
	v, e := jsonparser.GetString(msg.Payload, "type")
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
			resp = register(msg.ConnectID, msg.Payload)
		}
	case "login":
		{
			resp = login(msg.ConnectID, msg.Payload)
		}
	case "logout":
		{
			resp = logout(msg.ConnectID, msg.Payload)
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

func (m *Module) response(req *module.Message, payload map[string]interface{}) {
	msg := new(module.Message)
	msg.ConnectID = req.ConnectID
	msg.Sender = req.Sender
	msg.Recver = req.Recver
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

func (m *Module) handler(cmd string) {

}
