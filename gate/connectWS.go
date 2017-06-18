package gate

import (
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/hugohuang1111/woodcock/constants"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

type connectWS struct {
	socket  *websocket.Conn
	runFlag bool
	connID  uint64
	sendCh  chan []byte
}

func newConnectWS(socket *websocket.Conn) connect {
	c := new(connectWS)
	c.sendCh = make(chan []byte, 1024)
	c.socket = socket
	c.runFlag = true

	return c
}

func (c *connectWS) setID(id uint64) {
	c.connID = id
}

func (c *connectWS) ID() uint64 {
	return c.connID
}

func (c *connectWS) send(data []byte) {
	c.sendCh <- data
}

func (c *connectWS) onRecv(data []byte) {
	t, e := jsonparser.GetString(data, "type")
	if nil != e {
		glog.Error("wsConnection get type failed:", e)
		c.sendError(constants.ERROR_MSG_FORMAT_ERROR, "")
		c.runFlag = false
		return
	}

	s := strings.Split(t, ":")
	if 2 != len(s) {
		glog.Error("connect ws onRecv type error ", s)
		c.sendError(constants.ERROR_MSG_FORMAT_ERROR, t)
		c.runFlag = false
		return
	}
	name := s[0]
	m := new(module.Message)
	m.Recver = name
	m.Sender = module.MOD_GATE
	m.Type = module.MsgTypeClient
	m.Payload = make(map[string]interface{})
	m.Payload[module.PayloadKeyClientData] = data
	m.Payload[module.PayloadKeyConnectID] = c.connID

	router.Route(m)
}

func (c *connectWS) run() {
	clientConnect(c)
	go c.sendroutie()
	for c.runFlag {
		mt, message, err := c.socket.ReadMessage()
		if err != nil {
			glog.Warningf("Connect (%d) disconnect: %v", c.connID, err)
			c.runFlag = false
			break
		}

		switch mt {
		case websocket.TextMessage:
			{
				c.onRecv(message)
			}
		default:
		}
	}
	c.socket.Close()
	c.socket = nil
	clientDisconnect(c)
}

func (c *connectWS) close() {
	c.runFlag = false
}

func (c *connectWS) sendError(e int, t string) {
	resp := make(map[string]interface{})
	resp["version"] = 1
	resp["error"] = e
	resp["description"] = constants.ErrorMap[e]
	resp["type"] = t

	jsonString, err := json.Marshal(resp)
	if nil != err {
		glog.Error("marshal resp to string failed:", err)
		return
	}
	c.send(jsonString)
}

func (c *connectWS) sendroutie() {
	for c.runFlag {
		data := <-c.sendCh

		if nil == c.socket {
			glog.Error("send fail websocket is nil")
			return
		}
		c.socket.WriteMessage(websocket.TextMessage, data)
	}
}
