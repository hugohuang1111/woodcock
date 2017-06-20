package gate

import (
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const wsPort = ":8000"

var (
	wsUpgrader = websocket.Upgrader{} // use default options
	listener   net.Listener
)

//RunWS run websockt
func runServerWS() {
	listener, err := net.Listen("tcp", wsPort)
	if err != nil {
		glog.Fatalf("websocket listen failed:%v", err)
		return
	}
	glog.Infof("WebSocket listen on port %v", wsPort)

	httpServer := &http.Server{
		Handler:        http.HandlerFunc(ServeHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1024,
	}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	go httpServer.Serve(listener)
}

// StopWS -- stop websocker server
func stopWS() {
	listener.Close()
}

// ServeHandler -- http serve interface
func ServeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "not support", http.StatusNotFound)
		glog.Errorf("upgrade websocket failed:%v", err)
		return
	}

	c := newConnectWS(conn)
	c.run()
}
