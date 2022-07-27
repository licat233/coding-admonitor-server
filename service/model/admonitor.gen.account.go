package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type _AccountMgr struct {
	*_BaseMgr
}

// AccountMgr open func
func AccountMgr(db *gorm.DB) *_AccountMgr {
	if db == nil {
		panic(fmt.Errorf("AccountMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_AccountMgr{_BaseMgr: &_BaseMgr{DB: db.Table("account"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *_AccountMgr) GetTableName() string {
	return "account"
}

// Reset 重置gorm会话
func (obj *_AccountMgr) Reset() *_AccountMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_AccountMgr) Get() (result Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_AccountMgr) Gets() (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_AccountMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(Account{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取 表主键
func (obj *_AccountMgr) WithID(id int64) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithAid aid获取 账号id
func (obj *_AccountMgr) WithAid(aid string) Option {
	return optionFunc(func(o *options) { o.query["aid"] = aid })
}

// WithName name获取 账号名
func (obj *_AccountMgr) WithName(name string) Option {
	return optionFunc(func(o *options) { o.query["name"] = name })
}

// WithNickname nickname获取 账号昵称
func (obj *_AccountMgr) WithNickname(nickname string) Option {
	return optionFunc(func(o *options) { o.query["nickname"] = nickname })
}

// WithChannelID channel_id获取 归属渠道channel表主键id
func (obj *_AccountMgr) WithChannelID(channelID int64) Option {
	return optionFunc(func(o *options) { o.query["channel_id"] = channelID })
}

// WithStatus status获取 账号状态，0关闭,1开启
func (obj *_AccountMgr) WithStatus(status bool) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithCreateAt create_at获取 创建时间
func (obj *_AccountMgr) WithCreateAt(createAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["create_at"] = createAt })
}

// WithUpdateAt update_at获取 更新时间
func (obj *_AccountMgr) WithUpdateAt(updateAt time.Time) Option {
	return optionFunc(func(o *options) { o.query["update_at"] = updateAt })
}

// GetByOption 功能选项模式获取
func (obj *_AccountMgr) GetByOption(opts ...Option) (result Account, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_AccountMgr) GetByOptions(opts ...Option) (results []*Account, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where(options.query).Find(&results).Error

	return
}

// SelectPage 分页查询
func (obj *_AccountMgr) SelectPage(page IPage, opts ...Option) (resultPage IPage, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	resultPage = page
	results := make([]Account, 0)
	var count int64 // 统计总的记录数
	query := obj.DB.WithContext(obj.ctx).Model(Account{}).Where(options.query)
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

// GetFromID 通过id获取内容 表主键
func (obj *_AccountMgr) GetFromID(id int64) (result Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找 表主键
func (obj *_AccountMgr) GetBatchFromID(ids []int64) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromAid 通过aid获取内容 账号id
func (obj *_AccountMgr) GetFromAid(aid string) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`aid` = ?", aid).Find(&results).Error

	return
}

// GetBatchFromAid 批量查找 账号id
func (obj *_AccountMgr) GetBatchFromAid(aids []string) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`aid` IN (?)", aids).Find(&results).Error

	return
}

// GetFromName 通过name获取内容 账号名
func (obj *_AccountMgr) GetFromName(name string) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`name` = ?", name).Find(&results).Error

	return
}

// GetBatchFromName 批量查找 账号名
func (obj *_AccountMgr) GetBatchFromName(names []string) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`name` IN (?)", names).Find(&results).Error

	return
}

// GetFromNickname 通过nickname获取内容 账号昵称
func (obj *_AccountMgr) GetFromNickname(nickname string) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`nickname` = ?", nickname).Find(&results).Error

	return
}

// GetBatchFromNickname 批量查找 账号昵称
func (obj *_AccountMgr) GetBatchFromNickname(nicknames []string) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`nickname` IN (?)", nicknames).Find(&results).Error

	return
}

// GetFromChannelID 通过channel_id获取内容 归属渠道channel表主键id
func (obj *_AccountMgr) GetFromChannelID(channelID int64) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`channel_id` = ?", channelID).Find(&results).Error

	return
}

// GetBatchFromChannelID 批量查找 归属渠道channel表主键id
func (obj *_AccountMgr) GetBatchFromChannelID(channelIDs []int64) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`channel_id` IN (?)", channelIDs).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容 账号状态，0关闭,1开启
func (obj *_AccountMgr) GetFromStatus(status bool) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找 账号状态，0关闭,1开启
func (obj *_AccountMgr) GetBatchFromStatus(statuss []bool) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromCreateAt 通过create_at获取内容 创建时间
func (obj *_AccountMgr) GetFromCreateAt(createAt time.Time) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`create_at` = ?", createAt).Find(&results).Error

	return
}

// GetBatchFromCreateAt 批量查找 创建时间
func (obj *_AccountMgr) GetBatchFromCreateAt(createAts []time.Time) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`create_at` IN (?)", createAts).Find(&results).Error

	return
}

// GetFromUpdateAt 通过update_at获取内容 更新时间
func (obj *_AccountMgr) GetFromUpdateAt(updateAt time.Time) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`update_at` = ?", updateAt).Find(&results).Error

	return
}

// GetBatchFromUpdateAt 批量查找 更新时间
func (obj *_AccountMgr) GetBatchFromUpdateAt(updateAts []time.Time) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`update_at` IN (?)", updateAts).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_AccountMgr) FetchByPrimaryKey(id int64) (result Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// FetchIndexByAid  获取多个内容
func (obj *_AccountMgr) FetchIndexByAid(aid string, channelID int64, status bool) (results []*Account, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(Account{}).Where("`aid` = ? AND `channel_id` = ? AND `status` = ?", aid, channelID, status).Find(&results).Error

	return
}
