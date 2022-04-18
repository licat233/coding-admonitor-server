package spider

import (
	"admonitor-server/common/config"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_loginPage(t *testing.T) {
	config.Initialize()
	s := NewSpider()
	s.Login()
}

func Test_run(t *testing.T) {
	var stop = make(chan bool)
	var wg sync.WaitGroup
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("恢复函数接收到的信息：", err)
			stop <- true
		}
	}()
	wg.Add(1)
	go func() {
		//运行协程
		defer wg.Done()
		select {
		case <-stop:
			fmt.Println("stop通道接收讯息")
		}
		fmt.Println("goroutine已退出")
	}()
	logic()
	wg.Wait()
}

func logic() {
	var count int
	for {
		<-time.NewTimer(1 * time.Second).C
		fmt.Println("每隔1秒打印一次")
		count++
		if count == 5 {
			break
		}
	}
	panic(errors.New("服务器异常"))
}
