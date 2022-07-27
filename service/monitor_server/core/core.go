package core

import (
	"admonitor-server/common/emailx"
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/timer"
	"admonitor-server/service/monitor_server/core/logic"
	"fmt"
)

type CoreServer struct {
	CoreLogic *logic.CoreLogic
	Timer     timer.Timer
	StopCh    chan bool //通道
}

func NewCoreServer() *CoreServer {
	return &CoreServer{
		CoreLogic: logic.NewCoreLogic(),
		Timer:     timer.NewTimerTask(),
		StopCh:    make(chan bool),
	}
}

func (a *CoreServer) block() {
	<-a.StopCh
}

//Test 测试专用
func (a *CoreServer) Test() {
	defer func() {
		a.Stop()
	}()
	// defer recovery.RecoverHandler()
	global.LOG.Info("CoreServer启动中...")
	fmt.Println("CoreServer启动中...")
	//先运行一次，检查是否正常
	if !a.Once() {
		msg := "因为所有的Spider已关闭，CoreServer启动失败"
		global.LOG.Info(msg)
		fmt.Println(msg)
		return
	}
	//提示
	global.LOG.Info("CoreServer已启动")
	fmt.Println("CoreServer已启动")
	//锁住服务，防止服务关闭
	a.block()
}

func (a *CoreServer) Start() {
	defer func() {
		err := recover()
		if err != nil {
			if v, ok := err.(error); ok {
				a.Stop()
				global.LOG.Error(errorx.New("已经停止CoreServer服务", v))
				emailx.Send(fmt.Sprintf("CoreServer已停止:%s", v.Error()))
				fmt.Println(v.Error())
			} else {
				fmt.Println(err)
			}
		}
		a.Stop()
	}()
	// defer recovery.RecoverHandler()
	global.LOG.Info("CoreServer启动中...")
	//先运行一次，检查是否正常
	if !a.Once() {
		msg := "因为所有的Spider已关闭，CoreServer启动失败"
		global.LOG.Info(msg)
		return
	}
	//添加任务，cron v3版本默认不支持秒级定时，spec只能有 5 个field，分、时、天、月、星期几
	//每天n次的任务
	//每隔10分钟获取一次newest数据，时间间隔不要太短，爬虫太频繁访问会被封cookie
	a.Timer.AddTaskByFunc("CoreNewestData", "*/10 * * * *", a.CoreLogic.InsertNewestData)
	//每隔4分钟更新logic数据到数据库
	a.Timer.AddTaskByFunc("CurrentCoreLogicData", "*/4 * * * *", a.CoreLogic.UpdateCurrentLogicData)
	//每隔5分钟心跳一次，有些渠道需要，如果不需要，函数留空就行，心跳经常出现问题，最好关闭心跳
	// a.Timer.AddTaskByFunc("YesterdayPopinLogicData", "*/5 * * * *", a.CoreLogic.HeartbeatHandler)

	//每天一次的任务
	//因为存在缓存问题，所以在凌晨 00:02获取的昨天day数据可能不是最终的day数据
	a.Timer.AddTaskByFunc("CoreDayData", "2 0 * * *", a.CoreLogic.InsertDayData)
	//因为day数据可能有变化，所以在9:00的时候再更新一下昨天的day数据
	a.Timer.AddTaskByFunc("CoreDayData", "0 9 * * *", a.CoreLogic.InsertDayData)
	//每天早上 9:30，当天的创建node数据
	a.Timer.AddTaskByFunc("CoreNodeData", "30 9 * * *", a.CoreLogic.InsertNodeData)
	//每天早上 9:32，新建昨日的logic数据
	a.Timer.AddTaskByFunc("YesterdayCoreLogicData", "32 9 * * *", a.CoreLogic.InsertYesterdayLogicData)

	//提示
	global.LOG.Info("CoreServer已启动")
	fmt.Println("CoreServer已启动")
	//锁住服务，防止服务关闭
	a.block()
}

func (a *CoreServer) Stop() {
	a.Timer.Close()
	a.StopCh <- true
	close(a.StopCh)
	global.LOG.Info("CoreServer已停止")
	fmt.Println("CoreServer已停止")
	emailx.Send("AdMonitor服务已停止")
}

func (a *CoreServer) Restart() {
	global.LOG.Info("CoreServer重启中...")
	fmt.Println("CoreServer重启中...")
	a.Stop()
	a.Start()
	global.LOG.Info("CoreServer已重启")
	fmt.Println("CoreServer已重启")
}

func (a *CoreServer) Once() bool {
	a.CoreLogic.InsertDayData()
	a.CoreLogic.InsertNewestData()
	a.CoreLogic.InsertNodeData()
	a.CoreLogic.InsertYesterdayLogicData()
	a.CoreLogic.UpdateCurrentLogicData()
	return a.CoreLogic.AllSpiderStatus()
}
