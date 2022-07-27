package line_spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/service/model"
	"admonitor-server/service/monitor_server/core/spider"
	"errors"
	"time"
)

//GetStatistics 爬取批量子账号的statistics数据
func (s *LineSpider) GetStatistics(date time.Time, Type string, accounts []*model.Account) (statistics []*model.Statistics) {
	if s.Close {
		return nil
	}
	statisticsResp := s.Statistics(date)
	if statisticsResp == nil {
		global.LOG.Error(errorx.New("statisticsResp为nil", errors.New("Spider获取line广告数据失败")))
		return
	}
	var datas = statisticsResp.Datas
	for _, account := range accounts {
		for _, data := range datas {
			if data.Id == account.Aid {
				statistic := &data.Statistics
				status := false
				if data.DeliveryStatus == "ACTIVE" {
					status = true
				}
				statics := &model.Statistics{
					ChannelID:  account.ChannelID,
					AccountID:  account.ID,
					Cost:       statistic.Cost,
					Click:      statistic.Click,
					Imp:        statistic.Imp,
					Ctr:        statistic.Ctr,
					Cv:         statistic.Cv,
					Cvr:        statistic.Cvr,
					Cpa:        statistic.Cpa,
					Status:     status,
					Type:       Type,
					Currency:   statistic.Currency,
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

func (s *LineSpider) Info() *spider.Spider {
	return s.Spider
}

/** InitAccounts
 * 主要作用：
 * 网络获取该渠道下的所有子账户，并且保存到数据库
 */
func (s *LineSpider) InitAccounts() {
	if s.Close {
		panic(errors.New("spider is closed"))
	}
	respData := s.Statistics(time.Now().Local())
	if respData == nil {
		panic(errors.New("获取line广告数据失败"))
	}
	accountsData := respData.Datas
	accounts := []*model.Account{}
	for _, accountdata := range accountsData {
		if len(accountdata.Id) == 0 {
			continue
		}
		status := false
		if accountdata.DeliveryStatus == "ACTIVE" {
			status = true
		}
		account := &model.Account{
			Aid:       accountdata.Id,
			Name:      accountdata.Name,
			Nickname:  "未设置",
			ChannelID: s.Channel.ID,
			Status:    status,
			CreateAt:  time.Time{},
			UpdateAt:  time.Time{},
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
