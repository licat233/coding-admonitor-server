package popin_spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/utils"
	"admonitor-server/service/model"
	"admonitor-server/service/monitor_server/core/spider"
	"errors"
	"time"
)

//爬取批量子账号的statistics数据
func (s *PopinSpider) GetStatistics(date time.Time, Type string, accounts []*model.Account) (statistics []*model.Statistics) {
	if s.Close {
		return
	}
	statisticsResp := s.Statistics(date)
	if statisticsResp == nil {
		global.LOG.Error(errorx.New("statisticsResp为nil", errors.New("Spider获取popin广告数据失败")))
		return
	}
	var datas = statisticsResp.Data.Detail
	for _, account := range accounts {
		for _, data := range datas {
			if data.AccountID == account.Aid {
				status := false
				if data.Imp != "0" {
					status = true
				}
				statics := &model.Statistics{
					ChannelID:  account.ChannelID,
					AccountID:  account.ID,
					Cost:       utils.StringToFloat64(data.Cost),
					Click:      utils.StringToInt64(data.Click),
					Imp:        utils.StringToInt64(data.Imp),
					Ctr:        utils.StringToFloat64(data.Ctr),
					Cv:         utils.StringToInt64(data.Cv),
					Cvr:        utils.StringToFloat64(data.Cvr),
					Cpa:        utils.StringToFloat64(data.Cpa),
					Status:     status,
					Type:       Type,
					Currency:   data.Currency,
					RecordDate: date,
					CreateAt:   time.Time{},
					UpdateAt:   date,
					IDentifier: model.GenStatisticsIdentifier(account.ChannelID, account.ID, date, Type),
				}
				statistics = append(statistics, statics)
			}
		}
	}
	return
}

func (s *PopinSpider) Info() *spider.Spider {
	return s.Spider
}

/** InitAccounts
 * 主要作用：
 * 网络获取该渠道下的所有子账户，并且保存到数据库
 * 如果已经存在，则保存，如果不存在，则添加
 */
func (s *PopinSpider) InitAccounts() {
	if s.Close {
		panic(errors.New("spider is closed"))
	}
	respData := s.Statistics(time.Now().Local())
	if respData == nil {
		panic(errors.New("获取popin广告数据失败"))
	}
	accountsData := respData.Data.Detail
	accounts := []*model.Account{}
	for _, accountdata := range accountsData {
		if len(accountdata.AccountName) == 0 {
			continue
		}
		status := false
		if accountdata.Imp != "0" {
			status = true
		}
		account := &model.Account{
			Aid:       accountdata.AccountID,
			Name:      accountdata.AccountName,
			Nickname:  "popin账号",
			ChannelID: s.Channel.ID,
			Status:    status,
		}
		accounts = append(accounts, account)
	}
	m := model.NewAccountModel()
	err := m.ReplaceAll(accounts)
	if err != nil {
		global.LOG.Error(err)
		panic(err)
	}
}
