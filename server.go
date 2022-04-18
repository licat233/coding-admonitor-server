package main

import (
	"admonitor-server/common/config"
	"admonitor-server/service/line"
	"admonitor-server/service/popin"
	"fmt"
	"sync"
)

func init() {
	config.Initialize()
}

type Admonitor struct {
	WG     sync.WaitGroup //等待组
	StopCh chan bool      //通道
}

func NewAdmonitor() *Admonitor {
	return &Admonitor{}
}

func (a *Admonitor) BlockServer() {
	defer a.WG.Done()
	select {
	case <-a.StopCh:
		fmt.Println("Admonitor程序结束运行")
	}
}

func server() {
	admonitor := NewAdmonitor()
	admonitor.WG.Add(1)
	go func() {
		line.NewLineServer().Start()
	}()
	admonitor.WG.Add(1)
	go func() {
		popin.NewPopinServer().Start()
	}()
	admonitor.WG.Wait()
}
