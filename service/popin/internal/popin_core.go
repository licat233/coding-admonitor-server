package internal

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/service/model"
	"errors"
	"fmt"
)

/** NewPopin 新建line
 */
func NewPopinLogic() *PopinLogic {
	return &PopinLogic{
		Model:  model.NewPopinModel(),
		Spider: NewSpider(),
	}
}

//InsertNewestData 插入今天最新的数据到数据库
func (l *PopinLogic) InsertNewestData() {
	l.safeRun(func() {
		l.createPopinData(model.NEWEST, true)
	})
}

//InsertNodeData 插入今天节点数据到数据库
func (l *PopinLogic) InsertNodeData() {
	l.safeRun(func() {
		l.createPopinData(model.NODE, true)
	})
}

//InsertDayData 插入昨天的全天数据到数据库
func (l *PopinLogic) InsertDayData() {
	l.safeRun(func() {
		l.createPopinData(model.DAY, true)
	})
}

//InsertYesterdayLogicData 插入昨天的逻辑数据到数据库
func (l *PopinLogic) InsertYesterdayLogicData() {
	l.safeRun(l.insertYestdayLogicPopinData)
}

//UpdateCurrentLogicData 更新当前的逻辑数据到数据库
func (l *PopinLogic) UpdateCurrentLogicData() {
	l.safeRun(l.updateCurrentLogicPopinData)
}

//popin心跳
func (l *PopinLogic) HeartbeatRun() {
	l.safeRun(l.Spider.Heartbeat)
}

//safeRun 截断器
func (l *PopinLogic) safeRun(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			if v, ok := err.(error); ok {
				global.LOG.Error(errorx.New("已经暂停服务", v))
			} else {
				fmt.Println(err)
			}
		}
	}()
	if l.Spider.Close {
		panic(errors.New("spider已关闭"))
	}
	fn()
}
