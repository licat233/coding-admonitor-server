package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _LineMgr struct {
	*_BaseMgr
}

// LineMgr open func
func LineMgr(db *gorm.DB) *_LineMgr {
	if db == nil {
		panic(fmt.Errorf("LineMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_LineMgr{_BaseMgr: &_BaseMgr{DB: db.Table("line"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_LineMgr) GetTableName() string {
	return "line"
}

// Reset 重置gorm会话
func (obj *_LineMgr) Reset() *_LineMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_LineMgr) Get() (result Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_LineMgr) Gets() (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_LineMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Line{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取 主键
func (obj *_LineMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithAccount account获取 广告账号
func (obj *_LineMgr) WithAccount(account string) Option {
	return optionFunc(func(o *options) { o.query["account"] = account })
}

// WithCost cost获取 消耗额
func (obj *_LineMgr) WithCost(cost float64) Option {
	return optionFunc(func(o *options) { o.query["cost"] = cost })
}

// WithClick click获取 点击数
func (obj *_LineMgr) WithClick(click int64) Option {
	return optionFunc(func(o *options) { o.query["click"] = click })
}

// WithImp imp获取 展现量
func (obj *_LineMgr) WithImp(imp int64) Option {
	return optionFunc(func(o *options) { o.query["imp"] = imp })
}

// WithCtr ctr获取 点击率
func (obj *_LineMgr) WithCtr(ctr float64) Option {
	return optionFunc(func(o *options) { o.query["ctr"] = ctr })
}

// WithCv cv获取 转化数
func (obj *_LineMgr) WithCv(cv int64) Option {
	return optionFunc(func(o *options) { o.query["cv"] = cv })
}

// WithCvr cvr获取 转化率
func (obj *_LineMgr) WithCvr(cvr float64) Option {
	return optionFunc(func(o *options) { o.query["cvr"] = cvr })
}

// WithCpa cpa获取 转化成本
func (obj *_LineMgr) WithCpa(cpa float64) Option {
	return optionFunc(func(o *options) { o.query["cpa"] = cpa })
}

// WithStatus status获取 状态
func (obj *_LineMgr) WithStatus(status bool) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithType type获取 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_LineMgr) WithType(_type string) Option {
	return optionFunc(func(o *options) { o.query["type"] = _type })
}

// WithCurrency currency获取 币种
func (obj *_LineMgr) WithCurrency(currency string) Option {
	return optionFunc(func(o *options) { o.query["currency"] = currency })
}

// WithRecordDate record_date获取 记录日期
func (obj *_LineMgr) WithRecordDate(recordDate time.Time) Option {
	return optionFunc(func(o *options) { o.query["record_date"] = recordDate })
}

// WithCreateAt create_at获取 创建时间
func (obj *_LineMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新时间
func (obj *_LineMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_LineMgr) GetByOption(opts ...Option) (result Line, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_LineMgr) GetByOptions(opts ...Option) (results []*Line, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_LineMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Line, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Line{}).Where(options.query)
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
func (obj *_LineMgr) GetFromID(id int64) (result Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找 主键
func (obj *_LineMgr) GetBatchFromID(ids []int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromAccount 通过account获取内容 广告账号
func (obj *_LineMgr) GetFromAccount(account string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`account` = ?", account).Find(&results).Error

	return
}

// GetBatchFromAccount 批量查找 广告账号
func (obj *_LineMgr) GetBatchFromAccount(accounts []string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`account` IN (?)", accounts).Find(&results).Error

	return
}

// GetFromCost 通过cost获取内容 消耗额
func (obj *_LineMgr) GetFromCost(cost float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cost` = ?", cost).Find(&results).Error

	return
}

// GetBatchFromCost 批量查找 消耗额
func (obj *_LineMgr) GetBatchFromCost(costs []float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cost` IN (?)", costs).Find(&results).Error

	return
}

// GetFromClick 通过click获取内容 点击数
func (obj *_LineMgr) GetFromClick(click int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`click` = ?", click).Find(&results).Error

	return
}

// GetBatchFromClick 批量查找 点击数
func (obj *_LineMgr) GetBatchFromClick(clicks []int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`click` IN (?)", clicks).Find(&results).Error

	return
}

// GetFromImp 通过imp获取内容 展现量
func (obj *_LineMgr) GetFromImp(imp int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`imp` = ?", imp).Find(&results).Error

	return
}

// GetBatchFromImp 批量查找 展现量
func (obj *_LineMgr) GetBatchFromImp(imps []int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`imp` IN (?)", imps).Find(&results).Error

	return
}

// GetFromCtr 通过ctr获取内容 点击率
func (obj *_LineMgr) GetFromCtr(ctr float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`ctr` = ?", ctr).Find(&results).Error

	return
}

// GetBatchFromCtr 批量查找 点击率
func (obj *_LineMgr) GetBatchFromCtr(ctrs []float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`ctr` IN (?)", ctrs).Find(&results).Error

	return
}

// GetFromCv 通过cv获取内容 转化数
func (obj *_LineMgr) GetFromCv(cv int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cv` = ?", cv).Find(&results).Error

	return
}

// GetBatchFromCv 批量查找 转化数
func (obj *_LineMgr) GetBatchFromCv(cvs []int64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cv` IN (?)", cvs).Find(&results).Error

	return
}

// GetFromCvr 通过cvr获取内容 转化率
func (obj *_LineMgr) GetFromCvr(cvr float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cvr` = ?", cvr).Find(&results).Error

	return
}

// GetBatchFromCvr 批量查找 转化率
func (obj *_LineMgr) GetBatchFromCvr(cvrs []float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cvr` IN (?)", cvrs).Find(&results).Error

	return
}

// GetFromCpa 通过cpa获取内容 转化成本
func (obj *_LineMgr) GetFromCpa(cpa float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cpa` = ?", cpa).Find(&results).Error

	return
}

// GetBatchFromCpa 批量查找 转化成本
func (obj *_LineMgr) GetBatchFromCpa(cpas []float64) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`cpa` IN (?)", cpas).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容 状态
func (obj *_LineMgr) GetFromStatus(status bool) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找 状态
func (obj *_LineMgr) GetBatchFromStatus(statuss []bool) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromType 通过type获取内容 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_LineMgr) GetFromType(_type string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`type` = ?", _type).Find(&results).Error

	return
}

// GetBatchFromType 批量查找 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_LineMgr) GetBatchFromType(_types []string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`type` IN (?)", _types).Find(&results).Error

	return
}

// GetFromCurrency 通过currency获取内容 币种
func (obj *_LineMgr) GetFromCurrency(currency string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`currency` = ?", currency).Find(&results).Error

	return
}

// GetBatchFromCurrency 批量查找 币种
func (obj *_LineMgr) GetBatchFromCurrency(currencys []string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`currency` IN (?)", currencys).Find(&results).Error

	return
}

// GetFromRecordDate 通过record_date获取内容 记录日期
func (obj *_LineMgr) GetFromRecordDate(recordDate time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`record_date` = ?", recordDate).Find(&results).Error

	return
}

// GetBatchFromRecordDate 批量查找 记录日期
func (obj *_LineMgr) GetBatchFromRecordDate(recordDates []time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`record_date` IN (?)", recordDates).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_LineMgr) GetFromCreateAt(createAt time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_LineMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新时间
func (obj *_LineMgr) GetFromUpdateAt(updateAt time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新时间
func (obj *_LineMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_LineMgr) FetchByPrimaryKey(id int64) (result Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// FetchIndexByIndex1  获取多个内容
func (obj *_LineMgr) FetchIndexByIndex1(account string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`account` = ?", account).Find(&results).Error

	return
}

// FetchIndexByIndex2  获取多个内容
func (obj *_LineMgr) FetchIndexByIndex2(_type string) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`type` = ?", _type).Find(&results).Error

	return
}

// FetchIndexByIndex3  获取多个内容
func (obj *_LineMgr) FetchIndexByIndex3(recordDate time.Time) (results []*Line, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Line{}).Where("`record_date` = ?", recordDate).Find(&results).Error

	return
}
