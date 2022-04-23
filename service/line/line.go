package line

import (
	"admonitor-server/common/emailx"
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/timer"
	"admonitor-server/service/line/internal"
	"fmt"
)

type LineServer struct {
	LineLogic *internal.LineLogic
	Timer     timer.Timer
	StopCh    chan bool //通道
}

func NewLineServer() *LineServer {
	return &LineServer{
		LineLogic: internal.NewLineLogic(),
		Timer:     timer.NewTimerTask(),
	}
}

func (a *LineServer) block() {
	select {
	case <-a.StopCh:
		a.Stop()
	}
}

func (a *LineServer) Start() {
	defer func() {
		err := recover()
		if err != nil {
			if v, ok := err.(error); ok {
				global.LOG.Error(errorx.New("已经停止LineServer服务", v))
				emailx.Send(fmt.Sprintf("LineServer已停止:%s", v.Error()))
				fmt.Println(v.Error())
			} else {
				fmt.Println(err)
			}
			a.StopCh <- true
		}
	}()
	// defer recovery.RecoverHandler()
	global.LOG.Info("LineServer启动中...")
	fmt.Println("LineServer启动中...")
	//先运行一次，检查是否正常
	if !a.Once() {
		msg := "因为Spider已关闭，LineServer启动失败"
		global.LOG.Info(msg)
		fmt.Println(msg)
		return
	}
	//添加任务，cron v3版本默认不支持秒级定时，spec只能有 5 个field，分、时、天、月、星期几
	//每隔2分钟获取一次最新line数据
	a.Timer.AddTaskByFunc("LineNewestData", "*/2 * * * *", a.LineLogic.InsertNewestData)
	//每隔2分钟30秒更新逻辑数据到数据库
	a.Timer.AddTaskByFunc("CurrentLineLogicData", "*/2 * * * *", a.LineLogic.UpdateCurrentLogicData)
	//每天凌晨 00:02，获取昨日line数据
	a.Timer.AddTaskByFunc("LineDayData", "2 0 * * *", a.LineLogic.InsertDayData)
	//每天早上 9:30，获取该时间节点line数据
	a.Timer.AddTaskByFunc("LineNodeData", "30 9 * * *", a.LineLogic.InsertNodeData)
	//每天早上 9:32，新建昨日的逻辑line消耗数据
	a.Timer.AddTaskByFunc("YesterdayLineLogicData", "32 9 * * *", a.LineLogic.InsertYesterdayLogicData)
	//服务巡逻监控，如果spider已经关闭，则结束所有任务
	a.Timer.AddTaskByFunc("LineNewestData", "*/1 * * * *", a.monitorLineServer)
	//提示
	global.LOG.Info("LineServer已启动")
	fmt.Println("LineServer已启动")
	//锁住服务，防止服务关闭
	a.block()
}

func (a *LineServer) Stop() {
	a.Timer.Close()
	a.StopCh <- true
	close(a.StopCh)
	global.LOG.Info("LineServer已停止")
	fmt.Println("LineServer已停止")
}

func (a *LineServer) Restart() {
	global.LOG.Info("LineServer重启中...")
	fmt.Println("LineServer重启中...")
	a.Stop()
	a.Start()
	global.LOG.Info("LineServer已重启")
	fmt.Println("LineServer已重启")
}

func (a *LineServer) Once() bool {
	a.LineLogic.InsertDayData()
	a.LineLogic.InsertNewestData()
	a.LineLogic.UpdateCurrentLogicData()
	a.LineLogic.InsertNodeData()
	a.LineLogic.InsertYesterdayLogicData()
	return !a.LineLogic.Spider.Close
}

func (a *LineServer) monitorLineServer() {
	if a.LineLogic.Spider.Close {
		a.Stop()
		global.LOG.Info("Spider出现错误，LineServer已经停止运行")
	}
}
