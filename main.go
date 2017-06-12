package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/woodcock"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("WoodCock Start")
	woodcock.Run()
	glog.Info("WoodCock End")
}
