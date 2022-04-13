package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_0930(t *testing.T) {
	res := TargetTime(time.Now(), "00:02:00")
	fmt.Println(res)
}

func Test_time1(t *testing.T) {
	// InitTime()
	now := time.Now().Local()
	fmt.Println(now)
	target := TargetTime(now, "00:02:00")
	fmt.Println(target)
	res := now.Before(target)
	fmt.Println(res)
}

func Test_day(t *testing.T) {
	nowLogicDate := NowLogicDate()
	res := nowLogicDate.Day()
	fmt.Println(res)
}
