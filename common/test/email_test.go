package test

import (
	"admonitor-server/common/config"
	"admonitor-server/common/emailx"
	"testing"
)

func Test_emailx(t *testing.T) {
	config.Initialize()
	err := emailx.Send("测试内容")
	if err != nil {
		t.Error(err)
	}
}
