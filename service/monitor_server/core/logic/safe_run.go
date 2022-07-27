package logic

import (
	"admonitor-server/common/global"
	"admonitor-server/service/model"
	"admonitor-server/service/monitor_server/core/spider"
	"fmt"
	"time"
)

type Safe struct {
	ChannelSpider spider.SpiderMgr
	Spiders       map[int64]spider.SpiderMgr
}

func (core *CoreLogic) SafeRun(channelSpider spider.SpiderMgr) *Safe {
	return &Safe{
		ChannelSpider: channelSpider,
		Spiders:       core.Spiders,
	}
}

func (sf *Safe) Heartbeat() {
	defer func() {
		if err := recover(); err != nil {
			sf.RecoverHandler(err)
		}
	}()
	sf.ChannelSpider.Heartbeat()
}

//GetStatistics 爬取数据并且返回statistics数据
func (sf *Safe) GetStatistics(date time.Time, Type string, accounts []*model.Account) (res []*model.Statistics) {
	defer func() {
		if err := recover(); err != nil {
			res = nil
			sf.RecoverHandler(err)
		}
	}()
	return sf.ChannelSpider.GetStatistics(date, Type, accounts)
}

func (sf *Safe) RecoverHandler(info interface{}) {
	sf.ChannelSpider.Info().Close = true
	//移除出现panic的SpiderMgr
	// spiderId := sf.ChannelSpider.Info().Channel.ID
	//这里不使用遍历，使用映射性能会好一些
	// delete(sf.Spiders, spiderId)
	msg := fmt.Sprintf("%s爬虫出现了panic，请检查爬虫代码。%s", sf.ChannelSpider.Info().Channel.Name, info)
	global.LOG.Error(msg)
}
