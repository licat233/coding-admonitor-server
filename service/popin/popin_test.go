package popin

import (
	"admonitor-server/common/config"
	"admonitor-server/service/popin/internal"
	"testing"
)

func Test_popin(t *testing.T) {
	config.Initialize()
	popin := internal.NewPopinLogic()
	popin.InsertNewestData()
	popin.InsertNodeData()
	popin.InsertDayData()
	popin.InsertYesterdayLogicData()
	popin.UpdateCurrentLogicData()
	popin.HeartbeatRun()
}
