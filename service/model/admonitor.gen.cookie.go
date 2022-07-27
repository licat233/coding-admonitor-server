package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _CookieMgr struct {
	*_BaseMgr
}

// CookieMgr open func
func CookieMgr(db *gorm.DB) *_CookieMgr {
	if db == nil {
		panic(fmt.Errorf("CookieMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_CookieMgr{_BaseMgr: &_BaseMgr{DB: db.Table("cookie"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_CookieMgr) GetTableName() string {
	return "cookie"
}

// Reset 重置gorm会话
func (obj *_CookieMgr) Reset() *_CookieMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_CookieMgr) Get() (result Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_CookieMgr) Gets() (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_CookieMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Cookie{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取 主键
func (obj *_CookieMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithChannelID channel_id获取 cookie归属渠道channel表主键id
func (obj *_CookieMgr) WithChannelID(channelID int64) Option {
	return optionFunc(func(o *options) { o.query["channel_id"] = channelID })
}

// WithName name获取 设置的cookie_name
func (obj *_CookieMgr) WithName(name string) Option {
	return optionFunc(func(o *options) { o.query["name"] = name })
}

// WithValue value获取 设置的cookie_value
func (obj *_CookieMgr) WithValue(value string) Option {
	return optionFunc(func(o *options) { o.query["value"] = value })
}

// WithDomain domain获取 cookie的domain
func (obj *_CookieMgr) WithDomain(domain string) Option {
	return optionFunc(func(o *options) { o.query["domain"] = domain })
}

// WithPath path获取 cookie的path
func (obj *_CookieMgr) WithPath(path string) Option {
	return optionFunc(func(o *options) { o.query["path"] = path })
}

// WithCreateAt create_at获取 创建时间
func (obj *_CookieMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新时间
func (obj *_CookieMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_CookieMgr) GetByOption(opts ...Option) (result Cookie, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_CookieMgr) GetByOptions(opts ...Option) (results []*Cookie, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_CookieMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Cookie, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where(options.query)
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
func (obj *_CookieMgr) GetFromID(id int64) (result Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找 主键
func (obj *_CookieMgr) GetBatchFromID(ids []int64) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromChannelID 通过channel_id获取内容 cookie归属渠道channel表主键id
func (obj *_CookieMgr) GetFromChannelID(channelID int64) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`channel_id` = ?", channelID).Find(&results).Error

	return
}

// GetBatchFromChannelID 批量查找 cookie归属渠道channel表主键id
func (obj *_CookieMgr) GetBatchFromChannelID(channelIDs []int64) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`channel_id` IN (?)", channelIDs).Find(&results).Error

	return
}

// GetFromName 通过name获取内容 设置的cookie_name
func (obj *_CookieMgr) GetFromName(name string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// GetBatchFromName 批量查找 设置的cookie_name
func (obj *_CookieMgr) GetBatchFromName(names []string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromValue 通过value获取内容 设置的cookie_value
func (obj *_CookieMgr) GetFromValue(value string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`value` = ?", value).Find(&results).Error

	return
}

// GetBatchFromValue 批量查找 设置的cookie_value
func (obj *_CookieMgr) GetBatchFromValue(values []string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`value` IN (?)", values).Find(&results).Error

	return
}

// GetFromDomain 通过domain获取内容 cookie的domain
func (obj *_CookieMgr) GetFromDomain(domain string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`domain` = ?", domain).Find(&results).Error

	return
}

// GetBatchFromDomain 批量查找 cookie的domain
func (obj *_CookieMgr) GetBatchFromDomain(domains []string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`domain` IN (?)", domains).Find(&results).Error

	return
}

// GetFromPath 通过path获取内容 cookie的path
func (obj *_CookieMgr) GetFromPath(path string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`path` = ?", path).Find(&results).Error

	return
}

// GetBatchFromPath 批量查找 cookie的path
func (obj *_CookieMgr) GetBatchFromPath(paths []string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`path` IN (?)", paths).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_CookieMgr) GetFromCreateAt(createAt time.Time) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_CookieMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新时间
func (obj *_CookieMgr) GetFromUpdateAt(updateAt time.Time) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新时间
func (obj *_CookieMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_CookieMgr) FetchByPrimaryKey(id int64) (result Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// FetchIndexByIndex1  获取多个内容
func (obj *_CookieMgr) FetchIndexByIndex1(channelID int64) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`channel_id` = ?", channelID).Find(&results).Error

	return
}

// FetchIndexByIndex3  获取多个内容
func (obj *_CookieMgr) FetchIndexByIndex3(name string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// FetchIndexByIndex4  获取多个内容
func (obj *_CookieMgr) FetchIndexByIndex4(domain string) (results []*Cookie, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Cookie{}).Where("`domain` = ?", domain).Find(&results).Error

	return
}
