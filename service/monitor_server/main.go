package main

import (
	"admonitor-server/common/config"
	"admonitor-server/service/monitor_server/core"
	"admonitor-server/service/monitor_server/spider_list/line_spider"
	"admonitor-server/service/monitor_server/spider_list/popin_spider"
)

func main() {
	config.Initialize()
	c := core.NewCoreServer()
	//spider接口注册
	c.CoreLogic.SpiderListRegister("line", line_spider.NewSpiderMgr)
	c.CoreLogic.SpiderListRegister("popin", popin_spider.NewSpiderMgr)
	//初始化spiders
	c.CoreLogic.InitSpiders()
	//启动
	c.Start()
}
