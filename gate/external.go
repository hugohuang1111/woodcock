package gate

import (
	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/module"
)

//Module gate module
type Module struct {
	skelection module.Skelecton
}

//NewModule new gate module
func NewModule() *Module {
	m := new(Module)

	return m
}

//OnInit module init
func (m *Module) OnInit() {
	glog.Info("gate init")
	m.skelection.Run(m)
	runServerWS()
}

//OnDestroy module destroy
func (m *Module) OnDestroy() {
	stopWS()
	m.skelection.Stop()
}

//OnEvent module event
func (m *Module) OnEvent(msg *module.Message) {
	m.skelection.Add(msg)
}

//OnMsg module event
func (m *Module) OnMsg(msg *module.Message) {
	connID := module.GetConnectID(msg.Payload)
	clientData := module.GetClientData(msg.Payload)
	switch msg.Type {
	case module.MsgTypeClient:
		handlerClientMsg(connID, clientData)
	default:
		glog.Warning("gate unknow msg type:", msg.Type)
	}
}
