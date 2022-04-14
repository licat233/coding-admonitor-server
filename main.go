package main

import (
	"admonitor-server/common/config"
	"admonitor-server/service/line"
	"admonitor-server/service/popin"
)

func init() {
	config.Initialize()
}

func main() {
	Server()
}

func Server() {
	lineServer := line.NewLineServer()
	lineServer.Start()
	popinServer := popin.NewPopinServer()
	popinServer.Start()
}
