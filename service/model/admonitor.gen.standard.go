package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _StandardMgr struct {
	*_BaseMgr
}

// StandardMgr open func
func StandardMgr(db *gorm.DB) *_StandardMgr {
	if db == nil {
		panic(fmt.Errorf("StandardMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_StandardMgr{_BaseMgr: &_BaseMgr{DB: db.Table("standard"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_StandardMgr) GetTableName() string {
	return "standard"
}

// Reset 重置gorm会话
func (obj *_StandardMgr) Reset() *_StandardMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_StandardMgr) Get() (result Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_StandardMgr) Gets() (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_StandardMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Standard{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取 主键
func (obj *_StandardMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithAccountID account_id获取 归属渠道channel表主键id
func (obj *_StandardMgr) WithAccountID(accountID int64) Option {
	return optionFunc(func(o *options) { o.query["account_id"] = accountID })
}

// WithImp imp获取 展现
func (obj *_StandardMgr) WithImp(imp int64) Option {
	return optionFunc(func(o *options) { o.query["imp"] = imp })
}

// WithCost cost获取 消耗
func (obj *_StandardMgr) WithCost(cost float64) Option {
	return optionFunc(func(o *options) { o.query["cost"] = cost })
}

// WithClick click获取 点击
func (obj *_StandardMgr) WithClick(click int64) Option {
	return optionFunc(func(o *options) { o.query["click"] = click })
}

// WithCtr ctr获取 点击率
func (obj *_StandardMgr) WithCtr(ctr float64) Option {
	return optionFunc(func(o *options) { o.query["ctr"] = ctr })
}

// WithCv cv获取 转化
func (obj *_StandardMgr) WithCv(cv int64) Option {
	return optionFunc(func(o *options) { o.query["cv"] = cv })
}

// WithCvr cvr获取 转化率
func (obj *_StandardMgr) WithCvr(cvr float64) Option {
	return optionFunc(func(o *options) { o.query["cvr"] = cvr })
}

// WithCpa cpa获取 转化成本
func (obj *_StandardMgr) WithCpa(cpa float64) Option {
	return optionFunc(func(o *options) { o.query["cpa"] = cpa })
}

// WithStatus status获取 监控状态,0不监控，1监控
func (obj *_StandardMgr) WithStatus(status bool) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithCreateAt create_at获取 创建时间
func (obj *_StandardMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新日期
func (obj *_StandardMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_StandardMgr) GetByOption(opts ...Option) (result Standard, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_StandardMgr) GetByOptions(opts ...Option) (results []*Standard, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_StandardMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Standard, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Standard{}).Where(options.query)
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
func (obj *_StandardMgr) GetFromID(id int64) (result Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找 主键
func (obj *_StandardMgr) GetBatchFromID(ids []int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromAccountID 通过account_id获取内容 归属渠道channel表主键id
func (obj *_StandardMgr) GetFromAccountID(accountID int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`account_id` = ?", accountID).Find(&results).Error

	return
}

// GetBatchFromAccountID 批量查找 归属渠道channel表主键id
func (obj *_StandardMgr) GetBatchFromAccountID(accountIDs []int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`account_id` IN (?)", accountIDs).Find(&results).Error

	return
}

// GetFromImp 通过imp获取内容 展现
func (obj *_StandardMgr) GetFromImp(imp int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`imp` = ?", imp).Find(&results).Error

	return
}

// GetBatchFromImp 批量查找 展现
func (obj *_StandardMgr) GetBatchFromImp(imps []int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`imp` IN (?)", imps).Find(&results).Error

	return
}

// GetFromCost 通过cost获取内容 消耗
func (obj *_StandardMgr) GetFromCost(cost float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cost` = ?", cost).Find(&results).Error

	return
}

// GetBatchFromCost 批量查找 消耗
func (obj *_StandardMgr) GetBatchFromCost(costs []float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cost` IN (?)", costs).Find(&results).Error

	return
}

// GetFromClick 通过click获取内容 点击
func (obj *_StandardMgr) GetFromClick(click int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`click` = ?", click).Find(&results).Error

	return
}

// GetBatchFromClick 批量查找 点击
func (obj *_StandardMgr) GetBatchFromClick(clicks []int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`click` IN (?)", clicks).Find(&results).Error

	return
}

// GetFromCtr 通过ctr获取内容 点击率
func (obj *_StandardMgr) GetFromCtr(ctr float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`ctr` = ?", ctr).Find(&results).Error

	return
}

// GetBatchFromCtr 批量查找 点击率
func (obj *_StandardMgr) GetBatchFromCtr(ctrs []float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`ctr` IN (?)", ctrs).Find(&results).Error

	return
}

// GetFromCv 通过cv获取内容 转化
func (obj *_StandardMgr) GetFromCv(cv int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cv` = ?", cv).Find(&results).Error

	return
}

// GetBatchFromCv 批量查找 转化
func (obj *_StandardMgr) GetBatchFromCv(cvs []int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cv` IN (?)", cvs).Find(&results).Error

	return
}

// GetFromCvr 通过cvr获取内容 转化率
func (obj *_StandardMgr) GetFromCvr(cvr float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cvr` = ?", cvr).Find(&results).Error

	return
}

// GetBatchFromCvr 批量查找 转化率
func (obj *_StandardMgr) GetBatchFromCvr(cvrs []float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cvr` IN (?)", cvrs).Find(&results).Error

	return
}

// GetFromCpa 通过cpa获取内容 转化成本
func (obj *_StandardMgr) GetFromCpa(cpa float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cpa` = ?", cpa).Find(&results).Error

	return
}

// GetBatchFromCpa 批量查找 转化成本
func (obj *_StandardMgr) GetBatchFromCpa(cpas []float64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`cpa` IN (?)", cpas).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容 监控状态,0不监控，1监控
func (obj *_StandardMgr) GetFromStatus(status bool) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找 监控状态,0不监控，1监控
func (obj *_StandardMgr) GetBatchFromStatus(statuss []bool) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_StandardMgr) GetFromCreateAt(createAt time.Time) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_StandardMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新日期
func (obj *_StandardMgr) GetFromUpdateAt(updateAt time.Time) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新日期
func (obj *_StandardMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_StandardMgr) FetchByPrimaryKey(id int64) (result Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// FetchIndexByIndex2  获取多个内容
func (obj *_StandardMgr) FetchIndexByIndex2(accountID int64) (results []*Standard, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Standard{}).Where("`account_id` = ?", accountID).Find(&results).Error

	return
}
