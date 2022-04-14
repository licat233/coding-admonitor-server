package internal

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/utils"
	"admonitor-server/service/model"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//spiderLineData 爬取line数据
func (l *LineLogic) spiderLineData(date time.Time, Type string) (lines []*model.Line) {
	if Type == model.NODE {
		//规定只能在 09:00 - 10:00 时间段抓取
		if !utils.IsOn09to10(time.Now().Local()) {
			//如果不在09:00-10:00时间段，则不抓取，全部赋 0
			var nodeDate = utils.NowNodeDate()
			var accounts []*AccountConfig = AccountList() //由yaml文件提供
			for _, account := range accounts {
				if account.Status != true {
					continue
				}
				lines = append(lines, &model.Line{
					Account:    account.AccountID,
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
				})
			}
			return
		}
	}
	respData := l.Spider.Do(date)
	if respData == nil {
		global.LOG.Error(errorx.New("data为nil", errors.New("Spider获取line广告数据失败")))
		return
	}
	var datas = respData.Datas
	var accounts []*AccountConfig = AccountList() //由yaml文件提供
	for _, account := range accounts {
		if account.Status != true {
			continue
		}
		for _, data := range datas {
			if data.Id == account.AccountID {
				statistic := &data.Statistics
				var status bool
				if data.DeliveryStatus == "ACTIVE" {
					status = true
				}
				line := &model.Line{
					Account:    account.AccountID,
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
				}
				lines = append(lines, line)
			}
		}
	}
	return
}

/** createLineData 创建line的广告数据到数据库
 * @param Type Type string 操作的数据类型：model.NEWEST，model.NODE，model.DAY
 * @param replace type bool 如果存在，是否替换？
 */
func (l *LineLogic) createLineData(Type string, replace bool) {
	if l.Spider.Close {
		return
	}
	var date = time.Now().Local()
	if Type == model.DAY {
		date = date.AddDate(0, 0, -1)
	}
	lineRecords := l.spiderLineData(date, Type)
	for _, lineRecord := range lineRecords {
		l.InsertLineData(lineRecord, date, Type, replace)
	}
	return
}

//ReplaceData 插入数据到数据库，根据replace判断是否替换
func (l *LineLogic) InsertLineData(lineRecord *model.Line, date time.Time, Type string, replace bool) {
	if lineRecord == nil {
		global.LOG.Error("插入的lineRecord数据为空")
		return
	}
	res, err := l.Model.FindByDateAndType(lineRecord.Account, date, Type)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return
	}
	if err == nil && replace != true { //已经存在，且不更新
		return
	}
	if err == gorm.ErrRecordNotFound { //不存在，则插入
		if err = l.Model.InsertOne(lineRecord); err != nil {
			global.LOG.Error(errorx.New("GORM插入数据失败", err))
			return
		}
	} else {
		//已存在，则更新
		lineRecord.ID = res.ID
		if err = l.Model.UpdateByDateAndType(lineRecord, date, Type); err != nil {
			global.LOG.Error(errorx.New("GORM更新数据失败", err))
			return
		}
	}
	return
}

//InsertYestdayLogicLineData 创建昨天的逻辑数据到数据库
func (l *LineLogic) insertYestdayLogicLineData() {
	now := time.Now().Local()
	lists := AccountList()
	for _, info := range lists {
		if info.Status != true {
			continue
		}
		if now.Before(utils.Get0930(now)) {
			l.createLineData(model.DAY, true)
		}
		yesterday := now.AddDate(0, 0, -1)
		logicLine, err := l.LogicLineDataByDate(info.AccountID, yesterday)
		if err != nil {
			global.LOG.Error(errorx.New("获取昨天的line_logic数据失败", err))
			continue
		}
		if logicLine == nil {
			global.LOG.Error(errorx.New("创建昨天的line_logic数据失败", errors.New("logicLine为nil")))
			continue
		}
		l.InsertLineData(logicLine, yesterday, model.LOGIC, true)
	}
}

//updateCurrentLogicLineData 更新当前active的账户的line_logic数据到数据库
func (l *LineLogic) updateCurrentLogicLineData() {
	date := utils.NowLogicDate()
	lists := AccountList()
	for _, info := range lists {
		if info.Status != true {
			continue
		}
		logicLine, err := l.currentLogicLineData(info.AccountID)
		if err != nil {
			global.LOG.Error(errorx.New("获取当前line_logic数据失败", err))
			continue
		}
		if logicLine == nil {
			global.LOG.Error(errorx.New("更新当前line_logic数据失败", errors.New("logicLine为nil")))
			continue
		}
		l.InsertLineData(logicLine, date, model.LOGIC, true)
	}
}

//RangeDayCost 获取某账户范围日期内的全天总消耗
func (l *LineLogic) RangeDayCost(account string, start, end time.Time) (float64, error) {
	lines, err := l.Model.FindsByDateAndType(account, start, end, model.DAY)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var totalCost float64
	for _, line := range lines {
		totalCost += line.Cost
	}
	return totalCost, nil
}

/** TotalCost 获取时间范围内的总消费
 * @param start type time.Time 开始时间（包含）
 * @param end type time.Time 结束时间（包含）
 * 使用decimal库
 */
func (l *LineLogic) TotalCost(account string, start, end time.Time) (float64, error) {
	if start.After(end) {
		return 0, errors.New("开始时间大于结束时间")
	}
	now := time.Now().Local()
	if start.After(now) {
		return 0, nil
	}
	if end.After(now) {
		end = now
	}

	//全天消耗
	var dayCosts float64
	var err error
	//如果end等于now，无需查询daycost
	if end.Equal(now) {
		dayCosts = 0
	} else {
		dayCosts, err = l.RangeDayCost(account, start, end)
		if err != nil {
			return 0, err
		}
	}

	fmt.Println(dayCosts)

	//节点消耗
	//减去 start 当天的节点消耗
	return 0, nil
}

//CurrentLineData 获取某账户前最新logic消耗，recordDate已设置
func (l *LineLogic) currentLogicLineData(account string) (*model.Line, error) {
	logicDate := utils.NowLogicDate()
	now := time.Now().Local()
	if now.Before(utils.TargetTime(now, "00:02:00")) {
		return nil, errors.New("请等待到当天00:02:00之后再查询")
	}
	//获取当前最新消耗
	newestLine, err := l.Model.FindByDateAndType(account, now, model.NEWEST)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("数据库缺失%s的最新line数据", now.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}

	//获取logicDate的节点消耗
	nodeLine, err := l.Model.FindByDateAndType(account, logicDate, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有logicDate的节点数据，那就用logicDate最新的数据，作为line_logic数据返回
		newestLine.Type = model.LOGIC
		return newestLine, nil
		// return nil, fmt.Errorf("数据库缺失%s的节点数据", logicDate.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}

	//如果在 9:30 之前
	if now.Before(utils.Get0930(now)) {
		//获取logicDate的全天消耗
		dayLine, err := l.Model.FindByDateAndType(account, logicDate, model.DAY)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("数据库缺失%s的全天数据", logicDate.Format("2006-01-02"))
		}
		if err != nil {
			global.LOG.Error(errorx.New("GORM查询错误", err))
			return nil, err
		}
		return l.CalculateOnNodeBefore(dayLine, nodeLine, newestLine), nil
	}
	//如果在 9:30 之后
	return l.CalculateOnNodeAfter(newestLine, nodeLine), nil
}

//LogicLineData 获取某账户的某天logic消耗
func (l *LineLogic) LogicLineDataByDate(account string, date time.Time) (*model.Line, error) {
	//判断是否是今天
	if utils.NowLogicDate().Day() == date.Day() { //如果是今天
		return l.currentLogicLineData(account)
	}
	//指定日期的全天消耗
	dayLine, err := l.Model.FindByDateAndType(account, date, model.DAY)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("数据库缺失%s的全天数据", date.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}
	//指定日期的节点消耗
	nodeLine, err := l.Model.FindByDateAndType(account, date, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有date的节点数据，那就设为 0
		nodeLine = &model.Line{
			ID:         0,
			Account:    account,
			Cost:       0,
			Click:      0,
			Imp:        0,
			Ctr:        0,
			Cv:         0,
			Cvr:        0,
			Cpa:        0,
			Status:     false,
			Type:       model.NODE,
			Currency:   "",
			RecordDate: date,
			CreateAt:   utils.Get0930(date),
			UpdateAt:   date,
		}
		// return nil, fmt.Errorf("数据库缺失%s的节点数据", date.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}
	//指定日期的后一天节点消耗
	aftDate := utils.AfterDate(date)
	nodeLineAft, err := l.Model.FindByDateAndType(account, aftDate, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("数据库缺失%s后一天的节点数据", date.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}
	return l.CalculateOnNodeBefore(dayLine, nodeLine, nodeLineAft), nil
}

/** CalculateOnNodeBefore 计算节点前的line_logic数据 = 昨日全天消耗 - 昨日节点消耗 + (今天最新消耗||节点消耗)
 * @param: yerDayLine type *model.Line 前一天的全天数据
 * @param: yerNodeLine type *model.Line 前一天的节点数据
 * @param: nodeLine type *model.Line 今天最新消耗或者节点消耗
 */
func (l *LineLogic) CalculateOnNodeBefore(yerDayLine, yerNodeLine, nodeLine *model.Line) *model.Line {
	cost, ok := decimal.NewFromFloat(yerDayLine.Cost).Sub(decimal.NewFromFloat(yerNodeLine.Cost)).Add(decimal.NewFromFloat(nodeLine.Cost)).Round(2).Float64()
	if !ok {
		cost = yerDayLine.Cost - yerNodeLine.Cost + nodeLine.Cost
	}
	click := yerDayLine.Click - yerNodeLine.Click + nodeLine.Click
	imp := yerDayLine.Imp - yerNodeLine.Imp + nodeLine.Imp
	cv := yerDayLine.Cv - yerNodeLine.Cv + nodeLine.Cv
	ctr := utils.CtrCalculation(click, imp)
	cvr := utils.CvrCalculation(cv, click)
	cpa := utils.CpaCalculation(cost, cv)
	return &model.Line{
		Account:    yerNodeLine.Account,
		Cost:       cost,
		Click:      click,
		Imp:        imp,
		Ctr:        ctr,
		Cv:         cv,
		Cvr:        cvr,
		Cpa:        cpa,
		Status:     nodeLine.Status,
		Type:       model.LOGIC,
		Currency:   nodeLine.Currency,
		RecordDate: yerNodeLine.RecordDate,
	}
}

/** CalculateOnNodeAfter 计算节点后的line_logic数据 = 今日最新消耗 - 今日节点消耗
 * @param: nodeLine type *model.Line 当天的节点数据
 * @param: newestLine type *model.Line 当前的最新line数据
 */
func (l *LineLogic) CalculateOnNodeAfter(newestLine, nodeLine *model.Line) *model.Line {
	cost, ok := decimal.NewFromFloat(newestLine.Cost).Sub(decimal.NewFromFloat(nodeLine.Cost)).Round(2).Float64()
	if !ok {
		cost = newestLine.Cost - nodeLine.Cost
	}
	click := newestLine.Click - nodeLine.Click
	imp := newestLine.Imp - nodeLine.Imp
	cv := newestLine.Cv - nodeLine.Cv
	ctr := utils.CtrCalculation(click, imp)
	cvr := utils.CvrCalculation(cv, click)
	cpa := utils.CpaCalculation(cost, cv)
	return &model.Line{
		Account:    nodeLine.Account,
		Cost:       cost,
		Click:      click,
		Imp:        imp,
		Ctr:        ctr,
		Cv:         cv,
		Cvr:        cvr,
		Cpa:        cpa,
		Status:     nodeLine.Status,
		Type:       model.LOGIC,
		Currency:   nodeLine.Currency,
		RecordDate: nodeLine.RecordDate,
	}
}

/** AccountList 从yaml文件中读取account-list账号列表
 * @return: []*AccountConfig
 */
func AccountList() (res []*AccountConfig) {
	var n int
	if v, ok := viper.Get("line.account-list").([]interface{}); ok {
		n = len(v)
	} else {
		return
	}
	for i := 0; i < n; i++ {
		accountId := viper.GetString(fmt.Sprintf("line.account-list.%d.accountId", i))
		accountName := viper.GetString(fmt.Sprintf("line.account-list.%d.accountName", i))
		nickname := viper.GetString(fmt.Sprintf("line.account-list.%d.nickname", i))
		status := viper.GetBool(fmt.Sprintf("line.account-list.%d.status", i))
		res = append(res, &AccountConfig{
			AccountID:   accountId,
			AccountName: accountName,
			Nickname:    nickname,
			Status:      status,
		})
	}
	return
}
