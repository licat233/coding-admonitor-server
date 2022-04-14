package internal

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/service/model"
	"errors"
	"fmt"
)

type LineLogic struct {
	Model  *model.LineModel
	Spider *Spider
}

/** NewLine 新建line
 */
func NewLineLogic() *LineLogic {
	return &LineLogic{
		Model:  model.NewLineModel(),
		Spider: NewSpider(),
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

//safeRun 截断器
func (l *LineLogic) safeRun(fn func()) {
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
