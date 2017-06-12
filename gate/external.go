package gate

import (
	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/module"
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
	glog.Info("gate init")
	runServerWS()
}

//OnDestroy module destroy
func (m *Module) OnDestroy() {
	stopWS()
}

//OnEvent module event
func (m *Module) OnEvent(msg *module.Message) {
	switch msg.Type {
	case module.MOD_MSG_TYPE_CLIENT:
		{
			handlerClientMsg(msg.ConnectID, msg.Payload)
		}
	default:
		{
			glog.Warning("gate unknow msg type")
		}
	}
}
