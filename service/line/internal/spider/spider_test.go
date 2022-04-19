package spider

import (
	"admonitor-server/common/config"
	"testing"
	"time"
)

func Test_loginPage(t *testing.T) {
	config.Initialize()
	s := NewSpider()
	// s.loginPage()
	res := s.Statistics(time.Now())
	t.Log(res)
}
