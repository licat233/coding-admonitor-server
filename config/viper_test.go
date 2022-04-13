package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func Test_viper(t *testing.T) {
	LoadConfig()
	val := viper.Get("line.account-list")
	fmtType(val)
	if v, ok := val.([]interface{}); ok {
		fmt.Println(len(v))
	}
}

func fmtType(v interface{}) {
	res := reflect.TypeOf(v)
	fmt.Printf("类型：%v\n", res)
}
