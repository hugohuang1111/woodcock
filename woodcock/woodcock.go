package woodcock

import (
	"github.com/hugohuang1111/woodcock/gate"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
	"github.com/hugohuang1111/woodcock/user"
)

//Run run woodcock
func Run() {
	close := make(chan bool, 1)

	router.Start()
	module.Register(module.MOD_GATE, gate.NewModule())
	module.Register(module.MOD_USER, user.NewModule())
	module.Run()

	<-close
	router.End()
}
