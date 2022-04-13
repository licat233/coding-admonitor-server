package main

import (
	"admonitor-server/common/utils"
	"admonitor-server/config"
	"admonitor-server/service/line"
)

func init() {
	utils.InitTime()
	config.Initialize()
}

func main() {
	Server()
}

func Server() {
	lineServer := line.NewLineServer()
	lineServer.Start()
}
