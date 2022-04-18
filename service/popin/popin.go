package popin

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/timer"
	"admonitor-server/service/popin/internal"
	"fmt"
	"sync"
)

type PopinServer struct {
	PopinLogic *internal.PopinLogic
	Timer      timer.Timer
}

func NewPopinServer() *PopinServer {
	return &PopinServer{
		PopinLogic: internal.NewPopinLogic(),
		Timer:      timer.NewTimerTask(),
	}
}

//Run 服务启动程序
func Run() {
	var stopCh = make(chan bool)
	var wg sync.WaitGroup
	defer func() {
		err := recover()
		if err != nil {
			if v, ok := err.(error); ok {
				global.LOG.Error(errorx.New("已经停止PopinServer服务", v))
				fmt.Println(v.Error())
			} else {
				fmt.Println(err)
			}
			stopCh <- true
		}
	}()
	popinServer := NewPopinServer()
	popinServer.Start()
	//保持程序后台运行
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-stopCh:
			fmt.Println("PopinServer-Run程序结束运行")
		}
	}()

	wg.Wait()
}

func (a *PopinServer) Start() {
	global.LOG.Info("PopinServer启动中...")
	fmt.Println("PopinServer启动中...")
	// defer recovery.RecoverHandler()
	//先运行一次，检查是否正常
	if !a.Once() {
		msg := "因为Spider已关闭，PopinServer启动失败"
		global.LOG.Info(msg)
		fmt.Println(msg)
		return
	}
	//添加任务，cron v3版本默认不支持秒级定时，spec只能有 5 个field，分、时、天、月、星期几
	//每隔2分钟获取一次最新popin数据
	a.Timer.AddTaskByFunc("PopinNewestData", "*/2 * * * *", a.PopinLogic.InsertNewestData)
	//每隔2分钟30秒更新逻辑数据到数据库
	a.Timer.AddTaskByFunc("CurrentPopinLogicData", "*/2 * * * *", a.PopinLogic.UpdateCurrentLogicData)
	//每天凌晨 00:02，获取昨日popin数据
	a.Timer.AddTaskByFunc("PopinDayData", "2 0 * * *", a.PopinLogic.InsertDayData)
	//每天早上 9:30，获取该时间节点popin数据
	a.Timer.AddTaskByFunc("PopinNodeData", "30 9 * * *", a.PopinLogic.InsertNodeData)
	//每天早上 9:32，新建昨日的逻辑popin消耗数据
	a.Timer.AddTaskByFunc("YesterdayPopinLogicData", "32 9 * * *", a.PopinLogic.InsertYesterdayLogicData)
	//每隔 5 分钟心跳一次，更新cookie
	a.Timer.AddTaskByFunc("YesterdayPopinLogicData", "*/5 * * * *", a.PopinLogic.HeartbeatRun)
	//服务巡逻监控，如果spider已经关闭，则结束所有任务
	a.Timer.AddTaskByFunc("PopinNewestData", "0 */1 * * * *", a.monitorPopinServer)
	//提示
	global.LOG.Info("PopinServer已启动")
	fmt.Println("PopinServer已启动")
}

func (a *PopinServer) Stop() {
	a.Timer.Close()
	global.LOG.Info("PopinServer已停止")
	fmt.Println("PopinServer已停止")
}

func (a *PopinServer) Restart() {
	global.LOG.Info("PopinServer重启中...")
	fmt.Println("PopinServer重启中...")
	a.Stop()
	a.Start()
	global.LOG.Info("PopinServer已重启")
	fmt.Println("PopinServer已重启")
}

func (a *PopinServer) Once() bool {
	a.PopinLogic.InsertDayData()
	a.PopinLogic.InsertNewestData()
	a.PopinLogic.UpdateCurrentLogicData()
	a.PopinLogic.InsertNodeData()
	a.PopinLogic.InsertYesterdayLogicData()
	a.PopinLogic.HeartbeatRun()
	return !a.PopinLogic.Spider.Close
}

func (a *PopinServer) monitorPopinServer() {
	if a.PopinLogic.Spider.Close {
		a.Stop()
		global.LOG.Info("因为Spider已关闭，PopinServer已经停止运行")
	}
}
