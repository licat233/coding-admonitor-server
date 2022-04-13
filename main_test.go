package main

import (
	"monitor-data/common/utils"
	"monitor-data/config"
	"monitor-data/service/line/logic"
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
