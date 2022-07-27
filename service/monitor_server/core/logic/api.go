package logic

import (
	"admonitor-server/service/model"
	"errors"
)

//HeartbeatHandler 心跳处理
func (core *CoreLogic) HeartbeatHandler() {
	core.safeRun(func() {
		core.Heartbeat()
	})
}

//InsertNewestData 插入今天最新的数据到数据库
func (core *CoreLogic) InsertNewestData() {
	core.safeRun(func() {
		core.CreateCurrentStatistics(model.NEWEST)
	})
}

//InsertNodeData 插入今天节点数据到数据库
func (core *CoreLogic) InsertNodeData() {
	core.safeRun(func() {
		core.CreateCurrentStatistics(model.NODE)
	})
}

//InsertDayData 插入昨天的全天数据到数据库
func (core *CoreLogic) InsertDayData() {
	core.safeRun(func() {
		core.CreateCurrentStatistics(model.DAY)
	})
}

//InsertYesterdayLogicData 插入昨天的逻辑数据到数据库
func (core *CoreLogic) InsertYesterdayLogicData() {
	core.safeRun(core.insertYestdayLogicData)
}

//UpdateCurrentLogicData 更新当前的逻辑数据到数据库
func (core *CoreLogic) UpdateCurrentLogicData() {
	core.safeRun(core.updateCurrentLogicData)
}

//safeRun 截断器，这里的panic由LineServer恢复
func (core *CoreLogic) safeRun(fn func()) {
	if !core.AllSpiderStatus() {
		panic(errors.New("所有spider都已关闭"))
	}
	fn()
}

func (core *CoreLogic) AllSpiderStatus() bool {
	if len(core.Spiders) == 0 {
		panic("当前注册中心没有spider")
	}
	for _, spiderMgr := range core.Spiders {
		if !spiderMgr.Info().Close { //表示还有spider没有关闭，即等于false
			return true
		}
	}
	return false
}
