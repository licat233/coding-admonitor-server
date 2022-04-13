package config

import (
	"admonitor-server/common/global"
	"admonitor-server/common/gromx"
	"admonitor-server/common/logx"
	"admonitor-server/common/timer"
)

func Initialize() {
	LoadConfig()
	global.LOG = logx.NewLogger()
	global.DB = gromx.NewGromDB()
	global.TIMER = timer.NewTimerTask()
}
