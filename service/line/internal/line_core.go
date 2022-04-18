package internal

import (
	"admonitor-server/service/line/internal/spider"
	"admonitor-server/service/model"
	"errors"
)

/** NewLine 新建lineLogic
 */
func NewLineLogic() *LineLogic {
	return &LineLogic{
		Model:  model.NewLineModel(),
		Spider: spider.NewSpider(),
	}
}

//InsertNewestData 插入今天最新的数据到数据库
func (l *LineLogic) InsertNewestData() {
	l.safeRun(func() {
		l.createLineData(model.NEWEST, true)
	})
}

//InsertNodeData 插入今天节点数据到数据库
func (l *LineLogic) InsertNodeData() {
	l.safeRun(func() {
		l.createLineData(model.NODE, true)
	})
}

//InsertDayData 插入昨天的全天数据到数据库
func (l *LineLogic) InsertDayData() {
	l.safeRun(func() {
		l.createLineData(model.DAY, true)
	})
}

//InsertYesterdayLogicData 插入昨天的逻辑数据到数据库
func (l *LineLogic) InsertYesterdayLogicData() {
	l.safeRun(l.insertYestdayLogicLineData)
}

//UpdateCurrentLogicData 更新当前的逻辑数据到数据库
func (l *LineLogic) UpdateCurrentLogicData() {
	l.safeRun(l.updateCurrentLogicLineData)
}

//safeRun 截断器，这里的panic由LineServer恢复
func (l *LineLogic) safeRun(fn func()) {
	if l.Spider.Close {
		panic(errors.New("spider已关闭"))
	}
	fn()
}
