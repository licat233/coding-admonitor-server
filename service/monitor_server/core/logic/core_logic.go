package logic

import (
	"admonitor-server/common/emailx"
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/utils"
	"admonitor-server/service/model"
	"admonitor-server/service/monitor_server/core/spider"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type CoreLogic struct {
	ChannelModel    *model.ChannelModel
	AccountModel    *model.AccountModel
	StandardModel   *model.StandardModel
	StatisticsModel *model.StatisticsModel
	Spiders         map[int64]spider.SpiderMgr         //各个渠道的spider，必须用表id做映射
	SpiderMgrList   map[string]func() spider.SpiderMgr //存储各个渠道蜘蛛的新建方法,必须用channelName来做映射
}

func NewCoreLogic() *CoreLogic {
	return &CoreLogic{
		ChannelModel:    model.NewChannelModel(),
		AccountModel:    model.NewAccountModel(),
		StandardModel:   model.NewStandardModel(),
		StatisticsModel: model.NewStatisticsModel(),
		Spiders:         map[int64]spider.SpiderMgr{},
		SpiderMgrList:   map[string]func() spider.SpiderMgr{},
	}
}

//SpiderListRegister 渠道新建Spider的方法注册
//新编写的渠道扩展，自己通过这个函数注册进来
func (core *CoreLogic) SpiderListRegister(setChannelName string, newSpiderFunc func() spider.SpiderMgr) {
	if len(setChannelName) == 0 {
		panic("setChannelName is nil string")
	}
	core.SpiderMgrList[setChannelName] = newSpiderFunc
}

//InitSpiders 初始化所有需要爬取渠道的spider
//一定要在SpiderListRegister之后再运行
func (core *CoreLogic) InitSpiders() {
	//检查一下有无注册spider
	if len(core.SpiderMgrList) == 0 {
		return
	}
	//先从数据库加载需要爬取的渠道
	//TODO: 这里不应该只初始化状态为ture的渠道
	// channels := core.GetMonitorChannels()

	//TODO：应该获取core.SpiderMgrList里面的键值 === channel_name的渠道
	//①先获取所有的渠道
	allChannels, err := core.ChannelModel.FindAll()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//没有数据，返回
		panic(errors.New("数据库channel标中没有渠道数据"))
		// return
	}
	//②遍历挑选出符合条件的channels，并且给初始化
	for _, channel := range allChannels {
		for name, newSpider := range core.SpiderMgrList {
			if name == channel.Name { //符合条件的
				core.Spiders[channel.ID] = newSpider()
			}
		}
	}

	// if len(channels) == 0 {
	// 	msg := "当前数据库没有需要监控的channel，即status均为false"
	// 	global.LOG.Warn(msg)
	// 	panic(msg)
	// }

	// for _, channel := range channels {
	// 	//查找是否存在该渠道的spiderMgr
	// 	if newSpider, ok := core.SpiderMgrList[channel.Name]; ok {
	// 		//如果存在，那就新建该spiderMgr接口的spider，并且注册到Spiders中
	// 		core.Spiders[channel.ID] = newSpider()
	// 	}
	// 	//不存在不用管
	// }
}

//Heartbeat 批量心跳
//为了不发生阻塞，使用goroutine
func (core *CoreLogic) Heartbeat() {
	var wg sync.WaitGroup
	for _, channelSpider := range core.Spiders {
		wg.Add(1)
		go func(s spider.SpiderMgr) {
			defer wg.Done()
			core.SafeRun(s).Heartbeat()
		}(channelSpider)
	}
	wg.Wait()
}

/** CreateCurrentStatistics 调用爬虫创建当前对应type的广告数据到数据库，核心入口函数
 * @param Type Type string 操作的数据类型：model.NEWEST，model.NODE，model.DAY
 * 核心入口函数
 */
func (core *CoreLogic) CreateCurrentStatistics(Type string) {
	var date = time.Now().Local()
	if Type == model.DAY {
		date = date.AddDate(0, 0, -1)
	}
	statistics := core.spiderStatistics(date, Type)
	if len(statistics) == 0 {
		return
	}

	if len(statistics) != 0 {
		//注意：这里不能批量插入，过程需要经过一些处理，比如是否需要替换，是否需要更新，是否需要插入
		for _, statics := range statistics {
			core.InsertStatisticsData(statics)
		}
	}
}

/** insertYestdayLogicData 将昨天的logic数据插入到数据库中
 * 创建昨天的逻辑数据到数据库，一般于9:30之后调用
 * 核心入口函数
 */
func (core *CoreLogic) insertYestdayLogicData() {
	yesterday := time.Now().Local().AddDate(0, 0, -1)
	channels := core.GetMonitorChannels()
	accounts := core.GetMonitorAccounts(channels)
	for _, account := range accounts {
		//获取指定账户昨天的逻辑数据，通过一系列计算得到
		logicData, err := core.GenLogicStatistics(account, yesterday)
		if err != nil {
			msg := fmt.Sprintf("统计昨天%s的logic数据失败", account.Name)
			global.LOG.Error(errorx.New(msg, err))
			continue
		}
		core.InsertStatisticsData(logicData)
	}
}

/** insertTargetDateLogicData
 * 插入指定日期的logic数据
 */
func (core *CoreLogic) InsertTargetDateLogicData(targetDate time.Time) {
	channels := core.GetMonitorChannels()
	accounts := core.GetMonitorAccounts(channels)
	for _, account := range accounts {
		//获取指定账户昨天的逻辑数据，通过一系列计算得到
		logicData, err := core.GenLogicStatistics(account, targetDate)
		if err != nil {
			msg := fmt.Sprintf("统计昨天%s的logic数据失败", account.Name)
			global.LOG.Error(errorx.New(msg, err))
			continue
		}
		core.InsertStatisticsData(logicData)
	}
}

/** updateCurrentLogicData
 * 统计当前监控的所有账户line_logic数据到数据库
 * 核心入口函数
 * 注意，不要在 00:00 ~ 00:02 时间段使用
 */
func (core *CoreLogic) updateCurrentLogicData() {
	//判断当前是否处于 00:00 ~ 00:02 时间段
	if utils.InTimeRange(time.Now(), "00:00:00", "00:02:00") {
		//这个时间段不执行
		return
	}
	//从数据库获取需要监控的渠道
	channels := core.GetMonitorChannels()
	//从数据库获取需要监控的账户
	accounts := core.GetMonitorAccounts(channels)
	allLogicDatas := []*model.Statistics{}
	for _, account := range accounts {
		logicData, err := core.GenCurrentLogicData(account)
		if err != nil {
			global.LOG.Error(errorx.New("获取当前line_logic数据失败", err))
			continue
		}
		allLogicDatas = append(allLogicDatas, logicData)
		//监控
		core.MonitorStatistics(logicData)
	}
	//统一插入
	//必须使用replace插入，存在则更新，不存在则创建
	if err := core.StatisticsModel.ReplaceAll(allLogicDatas); err != nil {
		//数据库操作出错
		msg := errorx.New("mysql批量替换statistics数据失败", err)
		global.LOG.Error(msg)
		panic(errors.New(msg))
	}
}

/** spiderStatistics 核心函数
 * 爬取所有渠道下的账号Statistics数据，数据源于爬虫
 * 不接受爬取 logic 类型数据
 */
func (core *CoreLogic) spiderStatistics(date time.Time, Type string) (statistics []*model.Statistics) {
	//不接受爬取 logic 类型数据
	if Type == model.LOGIC {
		global.LOG.Error("不接受爬取 logic 类型数据")
		return
	}
	//从数据库中获取需要监控的渠道
	channels := core.GetMonitorChannels()
	if len(channels) == 0 {
		return
	}
	//从数据库中获取需要监控的账户
	accounts := core.GetMonitorAccounts(channels)
	if len(accounts) == 0 {
		return
	}
	//节点数据规定只能在 09:00 - 10:00 时间段抓取
	//如果是请求节点数据，且不在09:00-10:00时间段，则不抓取，全部赋 0
	if Type == model.NODE && !utils.IsOn09to10(time.Now().Local()) {
		var nodeDate = utils.NowNodeDate()
		for _, account := range accounts {
			statistics = append(statistics, &model.Statistics{
				ChannelID:  account.ChannelID,
				AccountID:  account.ID,
				Cost:       0,
				Click:      0,
				Imp:        0,
				Ctr:        0,
				Cv:         0,
				Cvr:        0,
				Cpa:        0,
				Status:     false,
				Type:       Type,
				Currency:   "TWD",
				RecordDate: nodeDate,
				UpdateAt:   date,
				IDentifier: model.GenStatisticsIdentifier(account.ChannelID, account.ID, nodeDate, Type),
			})
		}
		return
	}

	//爬取所有渠道网站得到的model.Statistics数据
	//遍历出所有需要监控的channel
	var wg sync.WaitGroup
	var mux sync.Mutex
	for _, channel := range channels {
		//查看是否有注册
		if spider, ok := core.Spiders[channel.ID]; ok {
			//查看该spider是否已经close
			if spider.Info().Close {
				//close则跳过
				continue
			}
			//开始爬取数据
			//使用goroutine，避免阻塞
			wg.Add(1)
			go func(date time.Time, Type string, channel_id int64, all_accounts []*model.Account) {
				defer wg.Done()
				account_list := FilterAccounts(channel_id, all_accounts)
				if len(account_list) == 0 {
					return
				}
				//查看数据库是否已存在day或者node数据，如果已存在，则不再爬取更新day数据
				if Type == model.DAY || Type == model.NODE {
					//找出不存在day数据的account
					var day_accounts []*model.Account
					for _, account := range account_list {
						_, err := core.StatisticsModel.FindByIdentifier(model.GenStatisticsIdentifier(channel_id, account.ID, date, Type))
						if errors.Is(err, gorm.ErrRecordNotFound) { //不存在
							day_accounts = append(day_accounts, account)
						}
					}
					account_list = day_accounts
				}
				if len(account_list) == 0 {
					return
				}

				statics := core.SafeRun(spider).GetStatistics(date, Type, account_list)
				mux.Lock()
				statistics = append(statistics, statics...)
				mux.Unlock()
			}(date, Type, channel.ID, accounts)
		}
	}
	wg.Wait()
	return
}

//CreateDayData 调用爬虫创建指定日期的day数据到数据库
func (core *CoreLogic) CreateDayData(account *model.Account, date time.Time) error {
	//开始爬取数据
	statistics, err := core.GetSingleStatistics(date, model.DAY, account)
	if err != nil {
		return err
	}
	//保存到数据库
	if err := core.StatisticsModel.Replace(statistics); err != nil {
		global.LOG.Error(errorx.New("mysql插入statistics数据失败", err))
		return err
	}
	return nil
}

/** InsertStatisticsData
 * 插入单条statistics数据到数据库，根据replace判断是否替换，内部只监控newest数据
 */
func (core *CoreLogic) InsertStatisticsData(statistics *model.Statistics) {
	if statistics == nil {
		global.LOG.Error("传入的statistics数据为空")
		return
	}
	//监控一下数据，看是否符合标准
	core.MonitorStatistics(statistics)
	//替换数据
	err := core.StatisticsModel.Replace(statistics)
	if err != nil {
		//数据库操作出错
		global.LOG.Error(errorx.New("mysql替换statistic数据失败", err))
		return
	}
}

//GenLogicStatistics 统计某账户的某天logic数据，通用函数，数据均源于数据库
func (core *CoreLogic) GenLogicStatistics(account *model.Account, date time.Time) (*model.Statistics, error) {
	//判断是否是今天
	if utils.NowLogicDate().Day() == date.Day() { //如果是今天，即 9:30-00:00
		return core.GenCurrentLogicData(account)
	}
	channel_id := account.ChannelID
	account_id := account.ID
	//指定日期的全天消耗
	dayData, err := core.StatisticsModel.FindByDateAndType(channel_id, account_id, date, model.DAY)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//缺少day数据，没法计算出logic数据
		//调用爬虫获取数据创建昨天的全天数据
		if err := core.CreateDayData(account, date); err != nil {
			return nil, err
		}
		//再查询一下指定日期全天数据
		dayData, err = core.StatisticsModel.FindByDateAndType(channel_id, account_id, date, model.DAY)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//还是不存在？不再调用爬虫了，直接报错返回
			return nil, fmt.Errorf("数据库缺失%s的全天数据", date.Format("2006-01-02"))
		} else if err != nil {
			global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
			return nil, err
		}
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
		return nil, err
	}
	//指定日期的节点消耗
	nodeData, err := core.StatisticsModel.FindByDateAndType(channel_id, account_id, date, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有节点数据，那就插入默认的节点数据0，因为时效性，无法再通过爬虫获取
		nodeData = &model.Statistics{
			ChannelID:  channel_id,
			AccountID:  account_id,
			Type:       model.NODE,
			Currency:   "TWD",
			RecordDate: date,
			CreateAt:   utils.Get0930(date),
			IDentifier: model.GenStatisticsIdentifier(channel_id, account_id, date, model.NODE),
		}
		if err := core.StatisticsModel.Replace(nodeData); err != nil {
			global.LOG.Error(errorx.New("mysql替换statistics数据失败", err))
			return nil, err
		}
	} else if err != nil { //这里一定要用else if，防止走这一步
		global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
		return nil, err
	}
	//指定日期的后一天的节点消耗
	aftDate := utils.AfterDate(date)
	nodeDataAft, err := core.StatisticsModel.FindByDateAndType(channel_id, account_id, aftDate, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有节点数据，那就插入默认的节点数据0，因为时效性，无法再通过爬虫获取
		nodeDataAft = &model.Statistics{
			ChannelID:  channel_id,
			AccountID:  account_id,
			Type:       model.NODE,
			Currency:   "TWD",
			RecordDate: aftDate,
			CreateAt:   utils.Get0930(aftDate),
			IDentifier: model.GenStatisticsIdentifier(channel_id, account_id, date, model.NODE),
		}
		if err := core.StatisticsModel.Replace(nodeDataAft); err != nil {
			global.LOG.Error(errorx.New("mysql替换statistics数据失败", err))
			return nil, err
		}
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
		return nil, err
	}
	return CalculateOnNodeBefore(dayData, nodeData, nodeDataAft), nil
}

//GenCurrentLogicData 统计某账户前最新logic消耗，源于mysql
func (core *CoreLogic) GenCurrentLogicData(account *model.Account) (*model.Statistics, error) {
	now := time.Now().Local()
	if now.Before(utils.TargetTime(now, "00:02:00")) {
		return nil, errors.New("请等待到当天00:02:00之后再查询")
	}
	channel_id := account.ChannelID
	account_id := account.ID
	logicDate := utils.NowLogicDate()
	//获取logicDate的node数据
	nodeData, err := core.StatisticsModel.FindByDateAndType(channel_id, account_id, logicDate, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有当前logicDate的node数据，那就插入默认的节点数据0，因为时效性，无法再通过爬虫获取
		nodeData = &model.Statistics{
			ChannelID:  channel_id,
			AccountID:  account_id,
			Type:       model.NODE,
			Currency:   "TWD",
			RecordDate: logicDate,
			CreateAt:   utils.Get0930(logicDate),
			IDentifier: model.GenStatisticsIdentifier(channel_id, account_id, logicDate, model.NODE),
		}
		if err := core.StatisticsModel.Replace(nodeData); err != nil {
			global.LOG.Error(errorx.New("mysql替换statistics数据失败", err))
			return nil, err
		}
	} else if err != nil { //这里一定要用else if，防止走这一步
		global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
		return nil, err
	}

	//获取当前newest数据
	newestData, err := core.StatisticsModel.FindByDateAndType(channel_id, account_id, now, model.NEWEST)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//数据库里面不存在，调用爬虫获取
		// core.safeRun(func() {
		// 	core.CreateCurrentStatistics(model.NEWEST)
		// })
		//爬取完成以后，回调一下，缺点：一旦spider陷入恐慌，这里会发生无限回调
		// return core.GenCurrentLogicData(account)
		return nil, fmt.Errorf("数据库缺失账户id为%d的%s最新line数据", account_id, now.Format("2006-01-02"))
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
		return nil, err
	}

	//如果在 9:30 之前
	if now.Before(utils.Get0930(now)) {
		//获取logicDate的day消耗
		dayData, err := core.StatisticsModel.FindByDateAndType(channel_id, account_id, logicDate, model.DAY)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//缺少logicDate的day数据，没法计算出logic数据
			//调用爬虫获取数据创建logicDate的day数据
			if err := core.CreateDayData(account, logicDate); err != nil {
				return nil, err
			}
			//再查询一下logicDate日期day数据
			dayData, err = core.StatisticsModel.FindByDateAndType(channel_id, account_id, logicDate, model.DAY)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//还是不存在？不再调用爬虫了，直接报错返回
				return nil, fmt.Errorf("数据库缺失%s的全天数据", logicDate.Format("2006-01-02"))
			} else if err != nil {
				global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
				return nil, err
			}
		} else if err != nil {
			global.LOG.Error(errorx.New("mysql查询statistics数据失败", err))
			return nil, err
		}
		return CalculateOnNodeBefore(dayData, nodeData, newestData), nil
	}
	//如果在 9:30 之后
	return CalculateOnNodeAfter(newestData, nodeData), nil
}

//MonitorStatistics 监控数据，如果不符合标准数据，则发送报警信息，使用协程发送邮件
//暂时只监控newest和部分logic数据
func (core *CoreLogic) MonitorStatistics(statistics *model.Statistics) {
	//数据为nil
	if statistics == nil {
		return
	}
	//没有消耗
	if statistics.Cost == 0 {
		return
	}

	Type := statistics.Type

	//暂时只监控newest和部分logic数据
	if Type != model.NEWEST && Type != model.LOGIC {
		return
	}

	//如果是logic数据，那要判断是不是当前的logic数据
	if Type == model.LOGIC {
		if !utils.InSameDay(statistics.RecordDate, utils.NowLogicDate()) {
			//如果不是今天的logic数据，则不报警
			return
		}
	}

	//获取数据标准
	standard, err := core.StandardModel.FindByAccountId(statistics.AccountID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，那就插入一条进去
		err := core.StandardModel.InsertOne(&model.Standard{
			AccountID: statistics.AccountID,
		})
		if err != nil {
			global.LOG.Error(errorx.New("mysql插入standard数据失败", err))
		}
		return
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询standard数据失败", err))
		return
	}
	//如果不需要监控，则不报警
	if !standard.Status {
		return
	}
	//查找渠道名
	channel, err := core.ChannelModel.FindOne(statistics.ChannelID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，则返回
		msg := fmt.Sprintf("mysql查询channel_id: %d 数据不存在", statistics.ChannelID)
		global.LOG.Warn(errorx.New(msg, err))
		return
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询channel数据失败", err))
		return
	}
	channel_name := channel.Name
	//查找账户名
	account, err := core.AccountModel.FindOne(statistics.AccountID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，则返回
		msg := fmt.Sprintf("mysql查询account_id: %d 数据不存在", statistics.AccountID)
		global.LOG.Warn(errorx.New(msg, err))
		return
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询account数据失败", err))
		return
	}
	account_name := account.Name

	var warnMsg []string
	if standard.Cost != 0 {
		if statistics.Cost > standard.Cost {
			e := fmt.Errorf("%s渠道%s账户的%s_Cost数据超过标准数据：%f > %f", channel_name, account_name, Type, statistics.Cost, standard.Cost)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if standard.Cpa != 0 {
		if statistics.Cpa > standard.Cpa {
			e := fmt.Errorf("%s渠道%s账户的%s_Cpa数据超过标准数据：%f > %f", channel_name, account_name, Type, statistics.Cpa, standard.Cpa)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if standard.Ctr != 0 {
		if statistics.Ctr < standard.Ctr {
			e := fmt.Errorf("%s渠道%s账户的%s_Ctr数据低于标准数据：%f < %f", channel_name, account_name, Type, statistics.Ctr, standard.Ctr)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if standard.Cvr != 0 {
		if statistics.Cvr < standard.Cvr {
			e := fmt.Errorf("%s渠道%s账户的%s_cvr数据低于标准数据：%f <br %f", channel_name, account_name, Type, statistics.Cvr, standard.Cvr)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if len(warnMsg) == 0 {
		return
	}
	//发送报警邮件
	msg := strings.Join(warnMsg, "\n</br>")
	b, _ := json.Marshal(statistics)
	logicData := fmt.Sprintf("当前%s渠道%s账户%s数据：%s", channel_name, account_name, Type, string(b))
	go emailx.Send(msg + logicData) //使用协程发送邮件
}
