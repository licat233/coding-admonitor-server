package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_date(t *testing.T) {
	InitTime()
	res, e := ParseDate("2022-04-28")
	if e != nil {
		t.Error(e)
	}
	t.Log(res)
}

func Test_InTimeRange(t *testing.T) {
	InitTime()
	ok := InTimeRange(time.Now(), "00:00:00", "11:00:00")
	fmt.Println(ok)
}
