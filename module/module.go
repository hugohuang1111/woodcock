package module

import (
	"github.com/golang/glog"
)

//Module module interface
type Module interface {
	OnInit()
	OnDestroy()
	OnEvent(msg *Message)
}

var (
	modMap map[string]Module
)

//Register mdoule register
func Register(name string, m Module) {
	if nil == modMap {
		modMap = make(map[string]Module)
	}
	modMap[name] = m
}

//Find found module
func Find(name string) Module {
	m, ok := modMap[name]
	if ok {
		return m
	}

	return nil
}

//Run run
func Run() {
	for _, m := range modMap {
		glog.Info("run module")
		m.OnInit()
	}
}

//Destroy destroy
func Destroy() {
	for _, m := range modMap {
		m.OnDestroy()
	}
}
