package timer

import (
	"fmt"
	"log"
	"testing"

	"github.com/robfig/cron/v3"
)

func Test_timer(t *testing.T) {
	timer := NewTimerTask()
	defer timer.Close()
	timer.AddTaskByFunc("test", "*/1 * * * *", echoSum)
	select {}
}

func Test_corn(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	defer c.Stop()
	id, err := c.AddFunc("*/2 * * * * *", echoSum)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已经添加echoSum方法到corn，id为", id)
	c.Start()
	select {}
}

func echoSum() {
	fmt.Println("执行一次")
}
