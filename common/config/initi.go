package config

import (
	"admonitor-server/common/emailx"
	"admonitor-server/common/global"
	"admonitor-server/common/gromx"
	"admonitor-server/common/logx"
	"admonitor-server/common/timer"
	"admonitor-server/common/utils"
	"admonitor-server/common/viperx"
)

func Initialize() {
	utils.InitTime()
	viperx.LoadConfig()
	global.LOG = logx.NewLogger()
	global.DB = gromx.NewGromDB()
	// global.REDIS = redisx.NewRedis() //暂时不需要
	global.TIMER = timer.NewTimerTask()
	global.EMAIL = emailx.NewEmail()
}
