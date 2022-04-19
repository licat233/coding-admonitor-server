package internal

import (
	"admonitor-server/common/emailx"
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/utils"
	"admonitor-server/service/model"
	"admonitor-server/service/popin/internal/spider"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type PopinLogic struct {
	Model         *model.PopinModel
	StandardModel *model.StandardModel
	Spider        *spider.Spider
}

/** NewPopin 新建popin
 */
func NewPopinLogic() *PopinLogic {
	return &PopinLogic{
		Model:         model.NewPopinModel(),
		StandardModel: model.NewStandardModel(),
		Spider:        spider.NewSpider(),
	}
}

//AccountConfig 账户信息
type AccountConfig struct {
	AccountID   string
	AccountName string
	Nickname    string
	Status      bool
}

//spiderPopinData 爬取popin数据
func (l *PopinLogic) spiderPopinData(date time.Time, Type string) (popins []*model.Popin) {
	if Type == model.NODE {
		//规定只能在 09:00 - 10:00 时间段抓取节点数据
		if !utils.IsOn09to10(time.Now().Local()) {
			//如果不在09:00-10:00时间段，则不抓取，全部赋 0
			var nodeDate = utils.NowNodeDate()
			var accounts []*AccountConfig = AccountList() //由yaml文件提供
			for _, account := range accounts {
				if account.Status != true {
					continue
				}
				popins = append(popins, &model.Popin{
					Account:    account.AccountName,
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
	respData := l.Spider.Statistics(date)
	if respData == nil { //spider已关闭引起的，则没有数据，就不进行任何操作
		global.LOG.Error(errorx.New("statistics为nil", errors.New("Spider获取popin广告数据失败")))
		return
	}
	var datas = respData.Data.Detail
	var accounts []*AccountConfig = AccountList() //由yaml文件提供
	for _, account := range accounts {
		if account.Status != true { //如果不监控该账户，则跳过
			continue
		}
		for _, data := range datas {
			if data.AccountID == account.AccountID {
				var status bool
				if data.Cost != "0" {
					status = true
				}
				popin := &model.Popin{
					Account:    account.AccountName,
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
				}
				popins = append(popins, popin)
			}
		}
	}
	return
}

/** createPopinData 创建popin的广告数据到数据库
 * @param Type Type string 操作的数据类型：model.NEWEST，model.NODE，model.DAY
 * @param replace type bool 如果存在，是否替换？
 */
func (l *PopinLogic) createPopinData(Type string, replace bool) {
	if l.Spider.Close {
		return
	}
	var date = time.Now().Local()
	if Type == model.DAY {
		date = date.AddDate(0, 0, -1)
	}
	popinRecords := l.spiderPopinData(date, Type)
	for _, popinRecord := range popinRecords {
		l.InsertPopinData(popinRecord, date, Type, replace)
	}
	return
}

//ReplaceData 插入数据到数据库，根据replace判断是否替换
func (l *PopinLogic) InsertPopinData(popinRecord *model.Popin, date time.Time, Type string, replace bool) {
	if popinRecord == nil {
		global.LOG.Error("插入的popinRecord数据为空")
		return
	}
	if Type == model.NEWEST {
		l.monitorHandler(popinRecord, model.NEWEST)
	}
	res, err := l.Model.FindByDateAndType(popinRecord.Account, date, Type)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return
	}
	if err == nil && replace != true { //已经存在，且不更新
		return
	}
	if err == gorm.ErrRecordNotFound { //不存在，则插入
		if err = l.Model.InsertOne(popinRecord); err != nil {
			global.LOG.Error(errorx.New("GORM插入数据失败", err))
			return
		}
	} else {
		//已存在，则更新
		popinRecord.ID = res.ID
		if err = l.Model.UpdateByDateAndType(popinRecord, date, Type); err != nil {
			global.LOG.Error(errorx.New("GORM更新数据失败", err))
			return
		}
	}
	return
}

//InsertYestdayLogicPopinData 创建昨天的逻辑数据到数据库
func (l *PopinLogic) insertYestdayLogicPopinData() {
	now := time.Now().Local()
	lists := AccountList()
	for _, info := range lists {
		if info.Status != true {
			continue
		}
		if now.Before(utils.Get0930(now)) {
			l.createPopinData(model.DAY, true)
		}
		yesterday := now.AddDate(0, 0, -1)
		logicPopin, err := l.LogicPopinDataByDate(info.AccountName, yesterday)
		if err != nil {
			global.LOG.Error(errorx.New("获取昨天的popin_logic数据失败", err))
			continue
		}
		if logicPopin == nil {
			global.LOG.Error(errorx.New("创建昨天的popin_logic数据失败", errors.New("logicPopin为nil")))
			continue
		}
		l.InsertPopinData(logicPopin, yesterday, model.LOGIC, true)
	}
}

//updateCurrentLogicPopinData 更新当前active的账户的popin_logic数据到数据库
func (l *PopinLogic) updateCurrentLogicPopinData() {
	date := utils.NowLogicDate()
	lists := AccountList()
	for _, info := range lists {
		if info.Status != true {
			continue
		}
		logicPopin, err := l.currentLogicPopinData(info.AccountName)
		if err != nil {
			global.LOG.Error(errorx.New("获取当前popin_logic数据失败", err))
			continue
		}
		if logicPopin == nil {
			global.LOG.Error(errorx.New("更新当前popin_logic数据失败", errors.New("logicPopin为nil")))
			continue
		}
		l.InsertPopinData(logicPopin, date, model.LOGIC, true)
		l.monitorHandler(logicPopin, model.LOGIC)
	}
}

//monitorHandler 监控数据，如果不符合标准数据，则发送报警信息
func (l *PopinLogic) monitorHandler(popin *model.Popin, Type string) {
	if popin == nil {
		return
	}
	if popin.Cost == 0 {
		return
	}
	typename := ""
	if Type == model.NEWEST {
		typename = "newest数据"
	} else if Type == model.NODE {
		typename = "node数据"
	}
	standard, err := l.StandardModel.FindStandard(l.Spider.Channel, popin.Account)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，那就插入一条进去
		l.StandardModel.InsertOne(&model.Standard{
			ID:       0,
			Channel:  l.Spider.Channel,
			Account:  popin.Account,
			Imp:      0,
			Cost:     0,
			Click:    0,
			Ctr:      0,
			Cv:       0,
			Cvr:      0,
			Cpa:      0,
			CreateAt: time.Time{},
			UpdateAt: time.Time{},
		})
		return
	} else if err != nil {
		global.LOG.Error(fmt.Errorf("获取%s渠道%s账户的标准数据失败：%s", l.Spider.Channel, popin.Account, err))
	}
	if standard == nil {
		return
	}
	var warnMsg []string
	if standard.Cost != 0 {
		if popin.Cost > standard.Cost {
			e := fmt.Errorf("%s渠道%s账户的%sCost数据超过标准数据：%f > %f", l.Spider.Channel, popin.Account, typename, popin.Cost, standard.Cost)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if standard.Cpa != 0 {
		if popin.Cpa > standard.Cpa {
			e := fmt.Errorf("%s渠道%s账户的%sCpa数据超过标准数据：%f > %f", l.Spider.Channel, popin.Account, typename, popin.Cpa, standard.Cpa)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if standard.Ctr != 0 {
		if popin.Ctr < standard.Ctr {
			e := fmt.Errorf("%s渠道%s账户的%sCtr数据低于标准数据：%f < %f", l.Spider.Channel, popin.Account, typename, popin.Ctr, standard.Ctr)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if standard.Cvr != 0 {
		if popin.Cvr < standard.Cvr {
			e := fmt.Errorf("%s渠道%s账户的%svr数据低于标准数据：%f <br %f", l.Spider.Channel, popin.Account, typename, popin.Cvr, standard.Cvr)
			global.LOG.Error(e)
			warnMsg = append(warnMsg, e.Error())
		}
	}
	if len(warnMsg) == 0 {
		return
	}
	//发送报警邮件
	msg := strings.Join(warnMsg, "\n</br>")
	b, _ := json.Marshal(popin)
	logicData := fmt.Sprintf("当前%s数据：%s", typename, string(b))
	go emailx.Send(msg + logicData)
	return
}

//RangeDayCost 获取某账户范围日期内的全天总消耗
func (l *PopinLogic) RangeDayCost(account string, start, end time.Time) (float64, error) {
	popins, err := l.Model.FindsByDateAndType(account, start, end, model.DAY)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var totalCost float64
	for _, popin := range popins {
		totalCost += popin.Cost
	}
	return totalCost, nil
}

/** TotalCost 获取时间范围内的总消费
 * @param start type time.Time 开始时间（包含）
 * @param end type time.Time 结束时间（包含）
 * 使用decimal库
 */
func (l *PopinLogic) TotalCost(account string, start, end time.Time) (float64, error) {
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

//CurrentPopinData 获取某账户前最新logic消耗，recordDate已设置
func (l *PopinLogic) currentLogicPopinData(account string) (*model.Popin, error) {
	logicDate := utils.NowLogicDate()
	now := time.Now().Local()
	if now.Before(utils.TargetTime(now, "00:02:00")) {
		return nil, errors.New("请等待到当天00:02:00之后再查询")
	}

	//获取当前最新消耗
	newestPopin, err := l.Model.FindByDateAndType(account, now, model.NEWEST)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("数据库缺失%s的最新popin数据", now.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}

	//获取logicDate的节点消耗
	nodePopin, err := l.Model.FindByDateAndType(account, logicDate, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有logicDate的节点数据，那就用logicDate最新的数据，作为popin_logic数据返回
		newestPopin.Type = model.LOGIC
		return newestPopin, nil
		// return nil, fmt.Errorf("数据库缺失%s的节点数据", logicDate.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}

	//如果在 9:30 之前
	if now.Before(utils.Get0930(now)) {
		//获取logicDate的全天消耗
		dayPopin, err := l.Model.FindByDateAndType(account, logicDate, model.DAY)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("数据库缺失%s的全天数据", logicDate.Format("2006-01-02"))
		}
		if err != nil {
			global.LOG.Error(errorx.New("GORM查询错误", err))
			return nil, err
		}
		return l.CalculateOnNodeBefore(dayPopin, nodePopin, newestPopin), nil
	}
	//如果在 9:30 之后
	return l.CalculateOnNodeAfter(newestPopin, nodePopin), nil
}

//LogicPopinData 获取某账户的某天logic消耗
func (l *PopinLogic) LogicPopinDataByDate(account string, date time.Time) (*model.Popin, error) {
	//判断是否是今天
	if utils.NowLogicDate().Day() == date.Day() { //如果是今天
		return l.currentLogicPopinData(account)
	}
	//指定日期的全天消耗
	dayPopin, err := l.Model.FindByDateAndType(account, date, model.DAY)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("数据库缺失%s的全天数据", date.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}
	//指定日期的节点消耗
	nodePopin, err := l.Model.FindByDateAndType(account, date, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//如果没有date的节点数据，那就设为 0
		nodePopin = &model.Popin{
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
	} else if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}
	//指定日期的后一天节点消耗
	aftDate := utils.AfterDate(date)
	nodePopinAft, err := l.Model.FindByDateAndType(account, aftDate, model.NODE)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("数据库缺失%s后一天的节点数据", date.Format("2006-01-02"))
	}
	if err != nil {
		global.LOG.Error(errorx.New("GORM查询错误", err))
		return nil, err
	}
	return l.CalculateOnNodeBefore(dayPopin, nodePopin, nodePopinAft), nil
}

/** CalculateOnNodeBefore 计算节点前的popin_logic数据 = 昨日全天消耗 - 昨日节点消耗 + (今天最新消耗||节点消耗)
 * @param: yerDayPopin type *model.Popin 前一天的全天数据
 * @param: yerNodePopin type *model.Popin 前一天的节点数据
 * @param: nodePopin type *model.Popin 今天最新消耗或者节点消耗
 */
func (l *PopinLogic) CalculateOnNodeBefore(yerDayPopin, yerNodePopin, nodePopin *model.Popin) *model.Popin {
	cost, ok := decimal.NewFromFloat(yerDayPopin.Cost).Sub(decimal.NewFromFloat(yerNodePopin.Cost)).Add(decimal.NewFromFloat(nodePopin.Cost)).Round(2).Float64()
	if !ok {
		cost = yerDayPopin.Cost - yerNodePopin.Cost + nodePopin.Cost
	}
	click := yerDayPopin.Click - yerNodePopin.Click + nodePopin.Click
	imp := yerDayPopin.Imp - yerNodePopin.Imp + nodePopin.Imp
	cv := yerDayPopin.Cv - yerNodePopin.Cv + nodePopin.Cv
	ctr := utils.CtrCalculation(click, imp)
	cvr := utils.CvrCalculation(cv, click)
	cpa := utils.CpaCalculation(cost, cv)
	var status bool
	if yerNodePopin.Cost != 0 {
		status = true
	}
	return &model.Popin{
		Account:    yerNodePopin.Account,
		Cost:       cost,
		Click:      click,
		Imp:        imp,
		Ctr:        ctr,
		Cv:         cv,
		Cvr:        cvr,
		Cpa:        cpa,
		Status:     status,
		Type:       model.LOGIC,
		Currency:   nodePopin.Currency,
		RecordDate: yerNodePopin.RecordDate,
	}
}

/** CalculateOnNodeAfter 计算节点后的popin_logic数据 = 今日最新消耗 - 今日节点消耗
 * @param: nodePopin type *model.Popin 当天的节点数据
 * @param: newestPopin type *model.Popin 当前的最新popin数据
 */
func (l *PopinLogic) CalculateOnNodeAfter(newestPopin, nodePopin *model.Popin) *model.Popin {
	cost, ok := decimal.NewFromFloat(newestPopin.Cost).Sub(decimal.NewFromFloat(nodePopin.Cost)).Round(2).Float64()
	if !ok {
		cost = newestPopin.Cost - nodePopin.Cost
	}
	click := newestPopin.Click - nodePopin.Click
	imp := newestPopin.Imp - nodePopin.Imp
	cv := newestPopin.Cv - nodePopin.Cv
	ctr := utils.CtrCalculation(click, imp)
	cvr := utils.CvrCalculation(cv, click)
	cpa := utils.CpaCalculation(cost, cv)
	var status bool
	if nodePopin.Cost != 0 {
		status = true
	}
	return &model.Popin{
		Account:    nodePopin.Account,
		Cost:       cost,
		Click:      click,
		Imp:        imp,
		Ctr:        ctr,
		Cv:         cv,
		Cvr:        cvr,
		Cpa:        cpa,
		Status:     status,
		Type:       model.LOGIC,
		Currency:   nodePopin.Currency,
		RecordDate: nodePopin.RecordDate,
	}
}

/** AccountList 从yaml文件中读取account-list账号列表
 * @return: []*AccountConfig
 */
func AccountList() (res []*AccountConfig) {
	var n int
	if v, ok := viper.Get("popin.account-list").([]interface{}); ok {
		n = len(v)
	} else {
		return
	}
	for i := 0; i < n; i++ {
		accountId := viper.GetString(fmt.Sprintf("popin.account-list.%d.id", i))
		accountName := viper.GetString(fmt.Sprintf("popin.account-list.%d.name", i))
		nickname := viper.GetString(fmt.Sprintf("popin.account-list.%d.nickname", i))
		status := viper.GetBool(fmt.Sprintf("popin.account-list.%d.status", i))
		res = append(res, &AccountConfig{
			AccountID:   accountId,
			AccountName: accountName,
			Nickname:    nickname,
			Status:      status,
		})
	}
	return
}
