package line

import (
	"admonitor-server/common/config"
	"admonitor-server/service/line/internal"
	"testing"
)

func Test_line(t *testing.T) {
	config.Initialize()
	line := internal.NewLineLogic()
	line.InsertNewestData()
	line.InsertNodeData()
	line.InsertDayData()
	line.InsertYesterdayLogicData()
	line.UpdateCurrentLogicData()
}
