package logic

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/utils"
	"admonitor-server/service/model"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

/** CalculateOnNodeBefore 计算节点前的logic数据 = 昨日全天消耗 - 昨日节点消耗 + (今天最新消耗||节点消耗)
 * @param: yerDayData type *model.Statistics 前一天的全天数据
 * @param: yerNodeData type *model.Statistics 前一天的节点数据
 * @param: nodeData type *model.Statistics 今天最新消耗或者节点消耗
 */
func CalculateOnNodeBefore(yerDayData, yerNodeData, nodeData *model.Statistics) *model.Statistics {
	cost, _ := decimal.NewFromFloat(yerDayData.Cost).Sub(decimal.NewFromFloat(yerNodeData.Cost)).Add(decimal.NewFromFloat(nodeData.Cost)).Round(2).Float64()
	click := yerDayData.Click - yerNodeData.Click + nodeData.Click
	imp := yerDayData.Imp - yerNodeData.Imp + nodeData.Imp
	cv := yerDayData.Cv - yerNodeData.Cv + nodeData.Cv
	ctr := utils.CtrCalculation(click, imp)
	cvr := utils.CvrCalculation(cv, click)
	cpa := utils.CpaCalculation(cost, cv)
	var status bool
	if yerNodeData.Cost != 0 {
		status = true
	}
	return &model.Statistics{
		ChannelID:  yerNodeData.ChannelID,
		AccountID:  yerNodeData.AccountID,
		Cost:       cost,
		Click:      click,
		Imp:        imp,
		Ctr:        ctr,
		Cv:         cv,
		Cvr:        cvr,
		Cpa:        cpa,
		Status:     status,
		Type:       model.LOGIC,
		Currency:   yerNodeData.Currency,
		RecordDate: yerNodeData.RecordDate,
		IDentifier: model.GenStatisticsIdentifier(yerNodeData.ChannelID, yerNodeData.AccountID, yerNodeData.RecordDate, model.LOGIC),
	}
}

/** CalculateOnNodeAfter 计算节点后的logic数据 = 今日最新消耗 - 今日节点消耗
 * @param: nodeData type *model.Statistics 当天的节点数据
 * @param: newestData type *model.Statistics 当前的最新line数据
 */
func CalculateOnNodeAfter(newestData, nodeData *model.Statistics) *model.Statistics {
	cost, _ := decimal.NewFromFloat(newestData.Cost).Sub(decimal.NewFromFloat(nodeData.Cost)).Round(2).Float64()
	click := newestData.Click - nodeData.Click
	imp := newestData.Imp - nodeData.Imp
	cv := newestData.Cv - nodeData.Cv
	ctr := utils.CtrCalculation(click, imp)
	cvr := utils.CvrCalculation(cv, click)
	cpa := utils.CpaCalculation(cost, cv)
	var status bool
	if nodeData.Cost != 0 {
		status = true
	}
	return &model.Statistics{
		ChannelID:  nodeData.ChannelID,
		AccountID:  nodeData.AccountID,
		Cost:       cost,
		Click:      click,
		Imp:        imp,
		Ctr:        ctr,
		Cv:         cv,
		Cvr:        cvr,
		Cpa:        cpa,
		Status:     status,
		Type:       model.LOGIC,
		Currency:   nodeData.Currency,
		RecordDate: nodeData.RecordDate,
		IDentifier: model.GenStatisticsIdentifier(nodeData.ChannelID, nodeData.AccountID, nodeData.RecordDate, model.LOGIC),
	}
}

/**
 * FilterAccounts 过滤出指定channel_id的accounts
 */
func FilterAccounts(channel_id int64, inaccounts []*model.Account) (accounts []*model.Account) {
	for _, account := range inaccounts {
		if account.ChannelID == channel_id {
			accounts = append(accounts, account)
		}
	}
	return
}

/** GetMonitorChannels
 * 获取需要监控的渠道
 */
func (core *CoreLogic) GetMonitorChannels() (channels []*model.Channel) {
	channels, err := core.ChannelModel.FindAllByStatus(true)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询channel数据出错", err))
		return
	}
	if len(channels) == 0 {
		return
	}
	return
}

/** GetMonitorAccounts
 * 获取需要监控的账号
 */
func (core *CoreLogic) GetMonitorAccounts(channels []*model.Channel) (accounts []*model.Account) {
	if len(channels) == 0 {
		return
	}
	var channel_ids []int64
	for _, channel := range channels {
		channel_ids = append(channel_ids, channel.ID)
	}
	accounts, err := core.AccountModel.FindAllByChannelsAndStatus(channel_ids, true)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询account数据出错", err))
		return
	}
	if len(accounts) == 0 {
		return
	}
	return
}

/**
 * 通过account_id获取channel_id
 */
func (core *CoreLogic) GetChannelIdByAccountId(account_id int64) int64 {
	account, err := core.AccountModel.FindOne(account_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，则返回
		msg := fmt.Sprintf("mysql查询account_id: %d 数据不存在", account_id)
		global.LOG.Warn(errorx.New(msg, err))
		return -1
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询account数据失败", err))
		return -1
	}
	return account.ChannelID
}

//GetSingleStatistics 爬取单个子账号的statistics数据
func (core *CoreLogic) GetSingleStatistics(date time.Time, Type string, account *model.Account) (*model.Statistics, error) {
	if account == nil {
		return nil, errors.New("传入的account为nil")
	}
	channel, err := core.ChannelModel.FindOne(account.ChannelID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("mysql查询channel_id: %d 数据不存在", account.ChannelID)
		global.LOG.Error(msg)
		return nil, errors.New(msg)
	} else if err != nil {
		msg := errorx.New("mysql查询channel数据失败", err)
		global.LOG.Error(msg)
		return nil, errors.New(msg)
	}
	if spider, ok := core.Spiders[channel.ID]; ok {
		statistics := spider.GetStatistics(date, Type, []*model.Account{account})
		if len(statistics) == 0 {
			err = fmt.Errorf("爬取%s渠道%s账户%s的statistics数据失败", channel.Name, account.Name, date.Format("2006-01-02"))
			global.LOG.Error(err)
			return nil, err
		}
		return statistics[0], nil
	}
	msg := fmt.Sprintf("【 %s 】渠道spider不存在", channel.Name)
	global.LOG.Error(msg)
	return nil, errors.New(msg)
}
