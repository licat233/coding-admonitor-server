package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _ChannelMgr struct {
	*_BaseMgr
}

// ChannelMgr open func
func ChannelMgr(db *gorm.DB) *_ChannelMgr {
	if db == nil {
		panic(fmt.Errorf("ChannelMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_ChannelMgr{_BaseMgr: &_BaseMgr{DB: db.Table("channel"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_ChannelMgr) GetTableName() string {
	return "channel"
}

// Reset 重置gorm会话
func (obj *_ChannelMgr) Reset() *_ChannelMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_ChannelMgr) Get() (result Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_ChannelMgr) Gets() (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_ChannelMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Channel{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取
func (obj *_ChannelMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithName name获取 渠道名称
func (obj *_ChannelMgr) WithName(name string) Option {
	return optionFunc(func(o *options) { o.query["name"] = name })
}

// WithEmail email获取 渠道登录邮箱
func (obj *_ChannelMgr) WithEmail(email string) Option {
	return optionFunc(func(o *options) { o.query["email"] = email })
}

// WithPassport passport获取 渠道登录账号
func (obj *_ChannelMgr) WithPassport(passport string) Option {
	return optionFunc(func(o *options) { o.query["passport"] = passport })
}

// WithPassword password获取 渠道登录密码
func (obj *_ChannelMgr) WithPassword(password string) Option {
	return optionFunc(func(o *options) { o.query["password"] = password })
}

// WithDomain domain获取 渠道cookie根域，例如: .line.biz
func (obj *_ChannelMgr) WithDomain(domain string) Option {
	return optionFunc(func(o *options) { o.query["domain"] = domain })
}

// WithStatus status获取 渠道状态,0关闭，1开启
func (obj *_ChannelMgr) WithStatus(status bool) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithCreateAt create_at获取 创建时间
func (obj *_ChannelMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新时间
func (obj *_ChannelMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_ChannelMgr) GetByOption(opts ...Option) (result Channel, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_ChannelMgr) GetByOptions(opts ...Option) (results []*Channel, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_ChannelMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Channel, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Channel{}).Where(options.query)
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

// GetFromID 通过id获取内容
func (obj *_ChannelMgr) GetFromID(id int64) (result Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_ChannelMgr) GetBatchFromID(ids []int64) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromName 通过name获取内容 渠道名称
func (obj *_ChannelMgr) GetFromName(name string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// GetBatchFromName 批量查找 渠道名称
func (obj *_ChannelMgr) GetBatchFromName(names []string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromEmail 通过email获取内容 渠道登录邮箱
func (obj *_ChannelMgr) GetFromEmail(email string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`email` = ?", email).Find(&results).Error

	return
}

// GetBatchFromEmail 批量查找 渠道登录邮箱
func (obj *_ChannelMgr) GetBatchFromEmail(emails []string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`email` IN (?)", emails).Find(&results).Error

	return
}

// GetFromPassport 通过passport获取内容 渠道登录账号
func (obj *_ChannelMgr) GetFromPassport(passport string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`passport` = ?", passport).Find(&results).Error

	return
}

// GetBatchFromPassport 批量查找 渠道登录账号
func (obj *_ChannelMgr) GetBatchFromPassport(passports []string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`passport` IN (?)", passports).Find(&results).Error

	return
}

// GetFromPassword 通过password获取内容 渠道登录密码
func (obj *_ChannelMgr) GetFromPassword(password string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`password` = ?", password).Find(&results).Error

	return
}

// GetBatchFromPassword 批量查找 渠道登录密码
func (obj *_ChannelMgr) GetBatchFromPassword(passwords []string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`password` IN (?)", passwords).Find(&results).Error

	return
}

// GetFromDomain 通过domain获取内容 渠道cookie根域，例如: .line.biz
func (obj *_ChannelMgr) GetFromDomain(domain string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`domain` = ?", domain).Find(&results).Error

	return
}

// GetBatchFromDomain 批量查找 渠道cookie根域，例如: .line.biz
func (obj *_ChannelMgr) GetBatchFromDomain(domains []string) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`domain` IN (?)", domains).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容 渠道状态,0关闭，1开启
func (obj *_ChannelMgr) GetFromStatus(status bool) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找 渠道状态,0关闭，1开启
func (obj *_ChannelMgr) GetBatchFromStatus(statuss []bool) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_ChannelMgr) GetFromCreateAt(createAt time.Time) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_ChannelMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新时间
func (obj *_ChannelMgr) GetFromUpdateAt(updateAt time.Time) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新时间
func (obj *_ChannelMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_ChannelMgr) FetchByPrimaryKey(id int64) (result Channel, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Channel{}).Where("`id` = ?", id).Find(&result).Error

	return
}
