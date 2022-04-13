package main

import (
	"monitor-data/common/utils"
	"monitor-data/config"
	"monitor-data/service/line"
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
