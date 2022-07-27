package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _StatisticsMgr struct {
	*_BaseMgr
}

// StatisticsMgr open func
func StatisticsMgr(db *gorm.DB) *_StatisticsMgr {
	if db == nil {
		panic(fmt.Errorf("StatisticsMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_StatisticsMgr{_BaseMgr: &_BaseMgr{DB: db.Table("statistics"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_StatisticsMgr) GetTableName() string {
	return "statistics"
}

// Reset 重置gorm会话
func (obj *_StatisticsMgr) Reset() *_StatisticsMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_StatisticsMgr) Get() (result Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_StatisticsMgr) Gets() (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_StatisticsMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Statistics{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取 主键
func (obj *_StatisticsMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithIDentifier identifier获取 唯一标识符，由channel_id、account_id、record_date、type用-符号连接而成，目的是为了防止重复插入
func (obj *_StatisticsMgr) WithIDentifier(identifier string) Option {
	return optionFunc(func(o *options) { o.query["identifier"] = identifier })
}

// WithChannelID channel_id获取 渠道id，指的是channel表的主键
func (obj *_StatisticsMgr) WithChannelID(channelID int64) Option {
	return optionFunc(func(o *options) { o.query["channel_id"] = channelID })
}

// WithAccountID account_id获取 账号id，指的是account表中的主键
func (obj *_StatisticsMgr) WithAccountID(accountID int64) Option {
	return optionFunc(func(o *options) { o.query["account_id"] = accountID })
}

// WithCost cost获取 消耗额
func (obj *_StatisticsMgr) WithCost(cost float64) Option {
	return optionFunc(func(o *options) { o.query["cost"] = cost })
}

// WithClick click获取 点击数
func (obj *_StatisticsMgr) WithClick(click int64) Option {
	return optionFunc(func(o *options) { o.query["click"] = click })
}

// WithImp imp获取 展现量
func (obj *_StatisticsMgr) WithImp(imp int64) Option {
	return optionFunc(func(o *options) { o.query["imp"] = imp })
}

// WithCtr ctr获取 点击率
func (obj *_StatisticsMgr) WithCtr(ctr float64) Option {
	return optionFunc(func(o *options) { o.query["ctr"] = ctr })
}

// WithCv cv获取 转化数
func (obj *_StatisticsMgr) WithCv(cv int64) Option {
	return optionFunc(func(o *options) { o.query["cv"] = cv })
}

// WithCvr cvr获取 转化率
func (obj *_StatisticsMgr) WithCvr(cvr float64) Option {
	return optionFunc(func(o *options) { o.query["cvr"] = cvr })
}

// WithCpa cpa获取 转化成本
func (obj *_StatisticsMgr) WithCpa(cpa float64) Option {
	return optionFunc(func(o *options) { o.query["cpa"] = cpa })
}

// WithStatus status获取 状态
func (obj *_StatisticsMgr) WithStatus(status bool) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithType type获取 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_StatisticsMgr) WithType(_type string) Option {
	return optionFunc(func(o *options) { o.query["type"] = _type })
}

// WithCurrency currency获取 币种
func (obj *_StatisticsMgr) WithCurrency(currency string) Option {
	return optionFunc(func(o *options) { o.query["currency"] = currency })
}

// WithRecordDate record_date获取 记录日期，09:30-09:30为一天
func (obj *_StatisticsMgr) WithRecordDate(recordDate time.Time) Option {
	return optionFunc(func(o *options) { o.query["record_date"] = recordDate })
}

// WithCreateAt create_at获取 创建时间
func (obj *_StatisticsMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新时间
func (obj *_StatisticsMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_StatisticsMgr) GetByOption(opts ...Option) (result Statistics, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_StatisticsMgr) GetByOptions(opts ...Option) (results []*Statistics, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_StatisticsMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Statistics, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where(options.query)
	query.Count(&count)
	resultPage.SetTotal(count)
	if len(page.GetOrederItemsString()) > 0 {
		query = query.Order(page.GetOrederItemsString())
	}
	err = query.Limit(int(page.GetSize())).Offset(int(page.Offset())).Find(&results).Error

	resultPage.SetRecords(results)
	return
}

//////////////////////////enume case ////////////////////////////////////////////

// GetFromID 通过id获取内容 主键
func (obj *_StatisticsMgr) GetFromID(id int64) (result Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找 主键
func (obj *_StatisticsMgr) GetBatchFromID(ids []int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromIDentifier 通过identifier获取内容 唯一标识符，由channel_id、account_id、record_date、type用-符号连接而成，目的是为了防止重复插入
func (obj *_StatisticsMgr) GetFromIDentifier(identifier string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`identifier` = ?", identifier).Find(&results).Error

	return
}

// GetBatchFromIDentifier 批量查找 唯一标识符，由channel_id、account_id、record_date、type用-符号连接而成，目的是为了防止重复插入
func (obj *_StatisticsMgr) GetBatchFromIDentifier(identifiers []string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`identifier` IN (?)", identifiers).Find(&results).Error

	return
}

// GetFromChannelID 通过channel_id获取内容 渠道id，指的是channel表的主键
func (obj *_StatisticsMgr) GetFromChannelID(channelID int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`channel_id` = ?", channelID).Find(&results).Error

	return
}

// GetBatchFromChannelID 批量查找 渠道id，指的是channel表的主键
func (obj *_StatisticsMgr) GetBatchFromChannelID(channelIDs []int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`channel_id` IN (?)", channelIDs).Find(&results).Error

	return
}

// GetFromAccountID 通过account_id获取内容 账号id，指的是account表中的主键
func (obj *_StatisticsMgr) GetFromAccountID(accountID int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`account_id` = ?", accountID).Find(&results).Error

	return
}

// GetBatchFromAccountID 批量查找 账号id，指的是account表中的主键
func (obj *_StatisticsMgr) GetBatchFromAccountID(accountIDs []int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`account_id` IN (?)", accountIDs).Find(&results).Error

	return
}

// GetFromCost 通过cost获取内容 消耗额
func (obj *_StatisticsMgr) GetFromCost(cost float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cost` = ?", cost).Find(&results).Error

	return
}

// GetBatchFromCost 批量查找 消耗额
func (obj *_StatisticsMgr) GetBatchFromCost(costs []float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cost` IN (?)", costs).Find(&results).Error

	return
}

// GetFromClick 通过click获取内容 点击数
func (obj *_StatisticsMgr) GetFromClick(click int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`click` = ?", click).Find(&results).Error

	return
}

// GetBatchFromClick 批量查找 点击数
func (obj *_StatisticsMgr) GetBatchFromClick(clicks []int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`click` IN (?)", clicks).Find(&results).Error

	return
}

// GetFromImp 通过imp获取内容 展现量
func (obj *_StatisticsMgr) GetFromImp(imp int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`imp` = ?", imp).Find(&results).Error

	return
}

// GetBatchFromImp 批量查找 展现量
func (obj *_StatisticsMgr) GetBatchFromImp(imps []int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`imp` IN (?)", imps).Find(&results).Error

	return
}

// GetFromCtr 通过ctr获取内容 点击率
func (obj *_StatisticsMgr) GetFromCtr(ctr float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`ctr` = ?", ctr).Find(&results).Error

	return
}

// GetBatchFromCtr 批量查找 点击率
func (obj *_StatisticsMgr) GetBatchFromCtr(ctrs []float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`ctr` IN (?)", ctrs).Find(&results).Error

	return
}

// GetFromCv 通过cv获取内容 转化数
func (obj *_StatisticsMgr) GetFromCv(cv int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cv` = ?", cv).Find(&results).Error

	return
}

// GetBatchFromCv 批量查找 转化数
func (obj *_StatisticsMgr) GetBatchFromCv(cvs []int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cv` IN (?)", cvs).Find(&results).Error

	return
}

// GetFromCvr 通过cvr获取内容 转化率
func (obj *_StatisticsMgr) GetFromCvr(cvr float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cvr` = ?", cvr).Find(&results).Error

	return
}

// GetBatchFromCvr 批量查找 转化率
func (obj *_StatisticsMgr) GetBatchFromCvr(cvrs []float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cvr` IN (?)", cvrs).Find(&results).Error

	return
}

// GetFromCpa 通过cpa获取内容 转化成本
func (obj *_StatisticsMgr) GetFromCpa(cpa float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cpa` = ?", cpa).Find(&results).Error

	return
}

// GetBatchFromCpa 批量查找 转化成本
func (obj *_StatisticsMgr) GetBatchFromCpa(cpas []float64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`cpa` IN (?)", cpas).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容 状态
func (obj *_StatisticsMgr) GetFromStatus(status bool) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找 状态
func (obj *_StatisticsMgr) GetBatchFromStatus(statuss []bool) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromType 通过type获取内容 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_StatisticsMgr) GetFromType(_type string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`type` = ?", _type).Find(&results).Error

	return
}

// GetBatchFromType 批量查找 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_StatisticsMgr) GetBatchFromType(_types []string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`type` IN (?)", _types).Find(&results).Error

	return
}

// GetFromCurrency 通过currency获取内容 币种
func (obj *_StatisticsMgr) GetFromCurrency(currency string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`currency` = ?", currency).Find(&results).Error

	return
}

// GetBatchFromCurrency 批量查找 币种
func (obj *_StatisticsMgr) GetBatchFromCurrency(currencys []string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`currency` IN (?)", currencys).Find(&results).Error

	return
}

// GetFromRecordDate 通过record_date获取内容 记录日期，09:30-09:30为一天
func (obj *_StatisticsMgr) GetFromRecordDate(recordDate time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`record_date` = ?", recordDate).Find(&results).Error

	return
}

// GetBatchFromRecordDate 批量查找 记录日期，09:30-09:30为一天
func (obj *_StatisticsMgr) GetBatchFromRecordDate(recordDates []time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`record_date` IN (?)", recordDates).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_StatisticsMgr) GetFromCreateAt(createAt time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_StatisticsMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新时间
func (obj *_StatisticsMgr) GetFromUpdateAt(updateAt time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新时间
func (obj *_StatisticsMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_StatisticsMgr) FetchByPrimaryKey(id int64) (result Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// FetchUniqueByIDentifier primary or index 获取唯一内容
func (obj *_StatisticsMgr) FetchUniqueByIDentifier(identifier string) (result Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`identifier` = ?", identifier).Find(&result).Error

	return
}

// FetchIndexByIDentifier2  获取多个内容
func (obj *_StatisticsMgr) FetchIndexByIDentifier2(identifier string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`identifier` = ?", identifier).Find(&results).Error

	return
}

// FetchIndexByIndex1  获取多个内容
func (obj *_StatisticsMgr) FetchIndexByIndex1(channelID int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`channel_id` = ?", channelID).Find(&results).Error

	return
}

// FetchIndexByIndex2  获取多个内容
func (obj *_StatisticsMgr) FetchIndexByIndex2(accountID int64) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`account_id` = ?", accountID).Find(&results).Error

	return
}

// FetchIndexByIndex3  获取多个内容
func (obj *_StatisticsMgr) FetchIndexByIndex3(_type string) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`type` = ?", _type).Find(&results).Error

	return
}

// FetchIndexByIndex4  获取多个内容
func (obj *_StatisticsMgr) FetchIndexByIndex4(recordDate time.Time) (results []*Statistics, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Statistics{}).Where("`record_date` = ?", recordDate).Find(&results).Error

	return
}
