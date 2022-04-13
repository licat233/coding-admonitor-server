package logic

import (
	"admonitor-server/config"
	"errors"
	"fmt"
	"testing"
)

func Test_viper(t *testing.T) {
	config.LoadConfig()
	res := AccountList()
	fmt.Println(res)
}

func Test_panic(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	cutoff()
	a()
}

func cutoff() {
	panic(errors.New("出错了"))
}

func a() {
	fmt.Println("aaaa")
}

func b() {
	fmt.Println("bbbb")
}
