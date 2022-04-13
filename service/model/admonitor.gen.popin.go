package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _PopinMgr struct {
	*_BaseMgr
}

// PopinMgr open func
func PopinMgr(db *gorm.DB) *_PopinMgr {
	if db == nil {
		panic(fmt.Errorf("PopinMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_PopinMgr{_BaseMgr: &_BaseMgr{DB: db.Table("popin"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_PopinMgr) GetTableName() string {
	return "popin"
}

// Reset 重置gorm会话
func (obj *_PopinMgr) Reset() *_PopinMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_PopinMgr) Get() (result Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_PopinMgr) Gets() (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_PopinMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Popin{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取 主键
func (obj *_PopinMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithAccount account获取 广告账号
func (obj *_PopinMgr) WithAccount(account string) Option {
	return optionFunc(func(o *options) { o.query["account"] = account })
}

// WithCost cost获取 消耗额
func (obj *_PopinMgr) WithCost(cost float64) Option {
	return optionFunc(func(o *options) { o.query["cost"] = cost })
}

// WithClick click获取 点击数
func (obj *_PopinMgr) WithClick(click int64) Option {
	return optionFunc(func(o *options) { o.query["click"] = click })
}

// WithImp imp获取 展现量
func (obj *_PopinMgr) WithImp(imp int64) Option {
	return optionFunc(func(o *options) { o.query["imp"] = imp })
}

// WithCtr ctr获取 点击率
func (obj *_PopinMgr) WithCtr(ctr float64) Option {
	return optionFunc(func(o *options) { o.query["ctr"] = ctr })
}

// WithCv cv获取 转化数
func (obj *_PopinMgr) WithCv(cv int64) Option {
	return optionFunc(func(o *options) { o.query["cv"] = cv })
}

// WithCvr cvr获取 转化率
func (obj *_PopinMgr) WithCvr(cvr float64) Option {
	return optionFunc(func(o *options) { o.query["cvr"] = cvr })
}

// WithCpa cpa获取 转化成本
func (obj *_PopinMgr) WithCpa(cpa float64) Option {
	return optionFunc(func(o *options) { o.query["cpa"] = cpa })
}

// WithStatus status获取 状态
func (obj *_PopinMgr) WithStatus(status bool) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithType type获取 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_PopinMgr) WithType(_type string) Option {
	return optionFunc(func(o *options) { o.query["type"] = _type })
}

// WithCurrency currency获取 币种
func (obj *_PopinMgr) WithCurrency(currency string) Option {
	return optionFunc(func(o *options) { o.query["currency"] = currency })
}

// WithRecordDate record_date获取 记录日期
func (obj *_PopinMgr) WithRecordDate(recordDate time.Time) Option {
	return optionFunc(func(o *options) { o.query["record_date"] = recordDate })
}

// WithCreateAt create_at获取 创建时间
func (obj *_PopinMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新时间
func (obj *_PopinMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_PopinMgr) GetByOption(opts ...Option) (result Popin, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_PopinMgr) GetByOptions(opts ...Option) (results []*Popin, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_PopinMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Popin, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Popin{}).Where(options.query)
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
func (obj *_PopinMgr) GetFromID(id int64) (result Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找 主键
func (obj *_PopinMgr) GetBatchFromID(ids []int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromAccount 通过account获取内容 广告账号
func (obj *_PopinMgr) GetFromAccount(account string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`account` = ?", account).Find(&results).Error

	return
}

// GetBatchFromAccount 批量查找 广告账号
func (obj *_PopinMgr) GetBatchFromAccount(accounts []string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`account` IN (?)", accounts).Find(&results).Error

	return
}

// GetFromCost 通过cost获取内容 消耗额
func (obj *_PopinMgr) GetFromCost(cost float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cost` = ?", cost).Find(&results).Error

	return
}

// GetBatchFromCost 批量查找 消耗额
func (obj *_PopinMgr) GetBatchFromCost(costs []float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cost` IN (?)", costs).Find(&results).Error

	return
}

// GetFromClick 通过click获取内容 点击数
func (obj *_PopinMgr) GetFromClick(click int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`click` = ?", click).Find(&results).Error

	return
}

// GetBatchFromClick 批量查找 点击数
func (obj *_PopinMgr) GetBatchFromClick(clicks []int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`click` IN (?)", clicks).Find(&results).Error

	return
}

// GetFromImp 通过imp获取内容 展现量
func (obj *_PopinMgr) GetFromImp(imp int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`imp` = ?", imp).Find(&results).Error

	return
}

// GetBatchFromImp 批量查找 展现量
func (obj *_PopinMgr) GetBatchFromImp(imps []int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`imp` IN (?)", imps).Find(&results).Error

	return
}

// GetFromCtr 通过ctr获取内容 点击率
func (obj *_PopinMgr) GetFromCtr(ctr float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`ctr` = ?", ctr).Find(&results).Error

	return
}

// GetBatchFromCtr 批量查找 点击率
func (obj *_PopinMgr) GetBatchFromCtr(ctrs []float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`ctr` IN (?)", ctrs).Find(&results).Error

	return
}

// GetFromCv 通过cv获取内容 转化数
func (obj *_PopinMgr) GetFromCv(cv int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cv` = ?", cv).Find(&results).Error

	return
}

// GetBatchFromCv 批量查找 转化数
func (obj *_PopinMgr) GetBatchFromCv(cvs []int64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cv` IN (?)", cvs).Find(&results).Error

	return
}

// GetFromCvr 通过cvr获取内容 转化率
func (obj *_PopinMgr) GetFromCvr(cvr float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cvr` = ?", cvr).Find(&results).Error

	return
}

// GetBatchFromCvr 批量查找 转化率
func (obj *_PopinMgr) GetBatchFromCvr(cvrs []float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cvr` IN (?)", cvrs).Find(&results).Error

	return
}

// GetFromCpa 通过cpa获取内容 转化成本
func (obj *_PopinMgr) GetFromCpa(cpa float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cpa` = ?", cpa).Find(&results).Error

	return
}

// GetBatchFromCpa 批量查找 转化成本
func (obj *_PopinMgr) GetBatchFromCpa(cpas []float64) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`cpa` IN (?)", cpas).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容 状态
func (obj *_PopinMgr) GetFromStatus(status bool) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找 状态
func (obj *_PopinMgr) GetBatchFromStatus(statuss []bool) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromType 通过type获取内容 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_PopinMgr) GetFromType(_type string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`type` = ?", _type).Find(&results).Error

	return
}

// GetBatchFromType 批量查找 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
func (obj *_PopinMgr) GetBatchFromType(_types []string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`type` IN (?)", _types).Find(&results).Error

	return
}

// GetFromCurrency 通过currency获取内容 币种
func (obj *_PopinMgr) GetFromCurrency(currency string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`currency` = ?", currency).Find(&results).Error

	return
}

// GetBatchFromCurrency 批量查找 币种
func (obj *_PopinMgr) GetBatchFromCurrency(currencys []string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`currency` IN (?)", currencys).Find(&results).Error

	return
}

// GetFromRecordDate 通过record_date获取内容 记录日期
func (obj *_PopinMgr) GetFromRecordDate(recordDate time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`record_date` = ?", recordDate).Find(&results).Error

	return
}

// GetBatchFromRecordDate 批量查找 记录日期
func (obj *_PopinMgr) GetBatchFromRecordDate(recordDates []time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`record_date` IN (?)", recordDates).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_PopinMgr) GetFromCreateAt(createAt time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_PopinMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新时间
func (obj *_PopinMgr) GetFromUpdateAt(updateAt time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新时间
func (obj *_PopinMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_PopinMgr) FetchByPrimaryKey(id int64) (result Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// FetchIndexByIndex1  获取多个内容
func (obj *_PopinMgr) FetchIndexByIndex1(account string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`account` = ?", account).Find(&results).Error

	return
}

// FetchIndexByIndex2  获取多个内容
func (obj *_PopinMgr) FetchIndexByIndex2(_type string) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`type` = ?", _type).Find(&results).Error

	return
}

// FetchIndexByIndex3  获取多个内容
func (obj *_PopinMgr) FetchIndexByIndex3(recordDate time.Time) (results []*Popin, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Popin{}).Where("`record_date` = ?", recordDate).Find(&results).Error

	return
}
