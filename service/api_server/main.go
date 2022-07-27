package main

import (
	"admonitor-server/common/config"
	"admonitor-server/service/api_server/logic"
)

func main() {
	config.Initialize()
	s := logic.NewApiServer()
	s.Start()
}
