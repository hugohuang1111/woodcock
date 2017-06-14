package router

import "github.com/hugohuang1111/woodcock/module"
import "github.com/golang/glog"

var (
	inChan  chan *module.Message
	outChan chan *module.Message
	runFlag bool
)

// Start router start
func Start() {
	inChan = make(chan *module.Message)
	outChan = make(chan *module.Message)
	runFlag = true
	go run()
}

// End router end
func End() {
	runFlag = false
}

//Route route message
func Route(msg *module.Message) {
	glog.Infof("router Route %s -> %s ", msg.Sender, msg.Recver)
	inChan <- msg
}

func run() {
	for runFlag {
		msg := <-inChan
		m := module.Find(msg.Recver)
		if nil == m {
			glog.Warning("router can't find module: ", msg.Recver)
			continue
		}
		go m.OnEvent(msg)
	}
}
