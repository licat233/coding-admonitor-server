package main

import (
	"admonitor-server/common/utils"
	"admonitor-server/config"
	"admonitor-server/service/line/logic"
	"testing"
)

func Test_LineServer(t *testing.T) {
	utils.InitTime()
	config.Initialize()
	line := logic.NewLineLogic()
	line.InsertNewestData()
	line.InsertNodeData()
	line.InsertDayData()
	line.InsertYesterdayLogicData()
	line.UpdateCurrentLogicData()
}
