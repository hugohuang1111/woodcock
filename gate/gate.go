package gate

import (
	"sync"

	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

var (
	connMap       map[uint64]connect
	connIDMutex   sync.Mutex
	connIDCounter uint64
)

func init() {
	connMap = make(map[uint64]connect)
}

func clientConnect(c connect) {
	connIDMutex.Lock()
	connIDCounter++
	c.setID(connIDCounter)
	connMap[connIDCounter] = c
	connIDMutex.Unlock()
}

func clientDisconnect(c connect) {
	delete(connMap, c.ID())

	m := new(module.Message)
	m.Recver = module.MOD_USER
	m.Sender = module.MOD_GATE
	m.Type = module.MsgTypeDisconnect
	m.Payload = make(map[string]interface{})
	m.Payload[module.PayloadKeyConnectID] = c.ID()
	router.Route(m)
}

func handlerClientMsg(connID uint64, payload []byte) {
	c, exist := connMap[connID]
	if !exist {
		glog.Warning("gate not find connect:", connID)
		return
	}
	c.send(payload)
}
