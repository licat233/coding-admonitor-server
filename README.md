# Admonitor-server
 广告数据监控系统

## 目前支持:
 * Line广告数据监控

## 组件说明：
* 配置管理 - viper
* 日志处理 - zap
* 数据库客户端 - gorm
* 定时任务 - corn
* 爬虫工具 - colly
* model工具 - gormt
* 金额计算 - decimal

## 进度：
* 定时任务
  ```golang
  func (a *LineServer) Start() {
	// defer recovery.RecoverHandler()
	//添加任务
	//每隔2分钟获取一次最新line数据
	a.Timer.AddTaskByFunc("LineNewestData", "0 */2 * * * *", a.LineLogic.InsertNewestData)
	//每隔2分钟30秒更新逻辑数据到数据库
	a.Timer.AddTaskByFunc("CurrentLineLogicData", "0 */2 * * * *", a.LineLogic.UpdateCurrentLogicData)
	//每天凌晨 00:02，获取昨日line数据
	a.Timer.AddTaskByFunc("LineDayData", "0 2 0 * * *", a.LineLogic.InsertDayData)
	//每天早上 9:30:30，获取该时间节点line数据
	a.Timer.AddTaskByFunc("LineNodeData", "30 30 9 * * *", a.LineLogic.InsertNodeData)
	//每天早上 9:32，新建昨日的逻辑line消耗数据
	a.Timer.AddTaskByFunc("YesterdayLineLogicData", "0 32 9 * * *", a.LineLogic.InsertYesterdayLogicData)
	//开始任务
	a.Timer.StartTask("LineNewestData")
	a.Timer.StartTask("CurrentLineLogicData")
	a.Timer.StartTask("LineDayData")
	a.Timer.StartTask("LineNodeData")
	a.Timer.StartTask("YesterdayLineLogicData")
	//提示
	fmt.Println("LineServer已启动")
  }
  ```
* 目前已经完成line数据的爬取【√】
* 准备开展popin数据的爬取 【?】
* 前端面板的开发【】
