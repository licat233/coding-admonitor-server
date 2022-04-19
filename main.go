package main

import (
	"admonitor-server/common/config"
	"admonitor-server/service/line"
	"admonitor-server/service/popin"
	"sync"
)

func init() {
	config.Initialize()
}

func main() {
	server()
}

func server() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		line.NewLineServer().Start()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		popin.NewPopinServer().Start()
	}()
	wg.Wait()
}
