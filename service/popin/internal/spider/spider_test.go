package spider

import (
	"admonitor-server/common/config"
	"testing"
	"time"
)

func Test_popin(t *testing.T) {
	config.Initialize()
	s := NewSpider()
	res := s.Statistics(time.Now())
	t.Log("结果:", res)
}
