package config

import (
	"monitor-data/common/global"
	"monitor-data/common/gromx"
	"monitor-data/common/logx"
	"monitor-data/common/timer"
)

func Initialize() {
	LoadConfig()
	global.LOG = logx.NewLogger()
	global.DB = gromx.NewGromDB()
	global.TIMER = timer.NewTimerTask()
}
