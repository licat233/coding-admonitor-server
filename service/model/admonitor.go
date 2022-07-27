package model

import (
	"time"
)

// Account 广告子账号表，并非登录账户
type Account struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"-"`      // 表主键
	Aid       string    `gorm:"column:aid" json:"aid"`              // 账号id
	Name      string    `gorm:"column:name" json:"name"`            // 账号名
	Nickname  string    `gorm:"column:nickname" json:"nickname"`    // 账号昵称
	ChannelID int64     `gorm:"column:channel_id" json:"channelId"` // 归属渠道channel表主键id
	Status    bool      `gorm:"column:status" json:"status"`        // 账号状态，0关闭,1开启
	CreateAt  time.Time `gorm:"column:create_at" json:"createAt"`   // 创建时间
	UpdateAt  time.Time `gorm:"column:update_at" json:"updateAt"`   // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Account) TableName() string {
	return "account"
}

// AccountColumns get sql column name.获取数据库列名
var AccountColumns = struct {
	ID        string
	Aid       string
	Name      string
	Nickname  string
	ChannelID string
	Status    string
	CreateAt  string
	UpdateAt  string
}{
	ID:        "id",
	Aid:       "aid",
	Name:      "name",
	Nickname:  "nickname",
	ChannelID: "channel_id",
	Status:    "status",
	CreateAt:  "create_at",
	UpdateAt:  "update_at",
}

// Channel 渠道表
type Channel struct {
	ID       int64     `gorm:"primaryKey;column:id" json:"-"`
	Name     string    `gorm:"column:name" json:"name"`          // 渠道名称
	Email    string    `gorm:"column:email" json:"email"`        // 渠道登录邮箱
	Passport string    `gorm:"column:passport" json:"passport"`  // 渠道登录账号
	Password string    `gorm:"column:password" json:"password"`  // 渠道登录密码
	Domain   string    `gorm:"column:domain" json:"domain"`      // 渠道cookie根域，例如: .line.biz
	Status   bool      `gorm:"column:status" json:"status"`      // 渠道状态,0关闭，1开启
	CreateAt time.Time `gorm:"column:create_at" json:"createAt"` // 创建时间
	UpdateAt time.Time `gorm:"column:update_at" json:"updateAt"` // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Channel) TableName() string {
	return "channel"
}

// ChannelColumns get sql column name.获取数据库列名
var ChannelColumns = struct {
	ID       string
	Name     string
	Email    string
	Passport string
	Password string
	Domain   string
	Status   string
	CreateAt string
	UpdateAt string
}{
	ID:       "id",
	Name:     "name",
	Email:    "email",
	Passport: "passport",
	Password: "password",
	Domain:   "domain",
	Status:   "status",
	CreateAt: "create_at",
	UpdateAt: "update_at",
}

// Cookie 爬虫cookies表
type Cookie struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"-"`      // 主键
	ChannelID int64     `gorm:"column:channel_id" json:"channelId"` // cookie归属渠道channel表主键id
	Name      string    `gorm:"column:name" json:"name"`            // 设置的cookie_name
	Value     string    `gorm:"column:value" json:"value"`          // 设置的cookie_value
	Domain    string    `gorm:"column:domain" json:"domain"`        // cookie的domain
	Path      string    `gorm:"column:path" json:"path"`            // cookie的path
	CreateAt  time.Time `gorm:"column:create_at" json:"createAt"`   // 创建时间
	UpdateAt  time.Time `gorm:"column:update_at" json:"updateAt"`   // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Cookie) TableName() string {
	return "cookie"
}

// CookieColumns get sql column name.获取数据库列名
var CookieColumns = struct {
	ID        string
	ChannelID string
	Name      string
	Value     string
	Domain    string
	Path      string
	CreateAt  string
	UpdateAt  string
}{
	ID:        "id",
	ChannelID: "channel_id",
	Name:      "name",
	Value:     "value",
	Domain:    "domain",
	Path:      "path",
	CreateAt:  "create_at",
	UpdateAt:  "update_at",
}

// Standard 广告数据规范
type Standard struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"-"`      // 主键
	AccountID int64     `gorm:"column:account_id" json:"accountId"` // 归属渠道channel表主键id
	Imp       int64     `gorm:"column:imp" json:"imp"`              // 展现
	Cost      float64   `gorm:"column:cost" json:"cost"`            // 消耗
	Click     int64     `gorm:"column:click" json:"click"`          // 点击
	Ctr       float64   `gorm:"column:ctr" json:"ctr"`              // 点击率
	Cv        int64     `gorm:"column:cv" json:"cv"`                // 转化
	Cvr       float64   `gorm:"column:cvr" json:"cvr"`              // 转化率
	Cpa       float64   `gorm:"column:cpa" json:"cpa"`              // 转化成本
	Status    bool      `gorm:"column:status" json:"status"`        // 监控状态,0不监控，1监控
	CreateAt  time.Time `gorm:"column:create_at" json:"createAt"`   // 创建时间
	UpdateAt  time.Time `gorm:"column:update_at" json:"updateAt"`   // 更新日期
}

// TableName get sql table name.获取数据库表名
func (m *Standard) TableName() string {
	return "standard"
}

// StandardColumns get sql column name.获取数据库列名
var StandardColumns = struct {
	ID        string
	AccountID string
	Imp       string
	Cost      string
	Click     string
	Ctr       string
	Cv        string
	Cvr       string
	Cpa       string
	Status    string
	CreateAt  string
	UpdateAt  string
}{
	ID:        "id",
	AccountID: "account_id",
	Imp:       "imp",
	Cost:      "cost",
	Click:     "click",
	Ctr:       "ctr",
	Cv:        "cv",
	Cvr:       "cvr",
	Cpa:       "cpa",
	Status:    "status",
	CreateAt:  "create_at",
	UpdateAt:  "update_at",
}

// Statistics 广告投放监控数据
type Statistics struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"-"`        // 主键
	IDentifier string    `gorm:"column:identifier" json:"identifier"`  // 唯一标识符，由channel_id、account_id、record_date、type用-符号连接而成，目的是为了防止重复插入
	ChannelID  int64     `gorm:"column:channel_id" json:"channelId"`   // 渠道id，指的是channel表的主键
	AccountID  int64     `gorm:"column:account_id" json:"accountId"`   // 账号id，指的是account表中的主键
	Cost       float64   `gorm:"column:cost" json:"cost"`              // 消耗额
	Click      int64     `gorm:"column:click" json:"click"`            // 点击数
	Imp        int64     `gorm:"column:imp" json:"imp"`                // 展现量
	Ctr        float64   `gorm:"column:ctr" json:"ctr"`                // 点击率
	Cv         int64     `gorm:"column:cv" json:"cv"`                  // 转化数
	Cvr        float64   `gorm:"column:cvr" json:"cvr"`                // 转化率
	Cpa        float64   `gorm:"column:cpa" json:"cpa"`                // 转化成本
	Status     bool      `gorm:"column:status" json:"status"`          // 状态
	Type       string    `gorm:"column:type" json:"type"`              // 数据类型：newest最新数据，node节点数据，day全天数据，logic逻辑数据
	Currency   string    `gorm:"column:currency" json:"currency"`      // 币种
	RecordDate time.Time `gorm:"column:record_date" json:"recordDate"` // 记录日期，09:30-09:30为一天
	CreateAt   time.Time `gorm:"column:create_at" json:"createAt"`     // 创建时间
	UpdateAt   time.Time `gorm:"column:update_at" json:"updateAt"`     // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Statistics) TableName() string {
	return "statistics"
}

// StatisticsColumns get sql column name.获取数据库列名
var StatisticsColumns = struct {
	ID         string
	IDentifier string
	ChannelID  string
	AccountID  string
	Cost       string
	Click      string
	Imp        string
	Ctr        string
	Cv         string
	Cvr        string
	Cpa        string
	Status     string
	Type       string
	Currency   string
	RecordDate string
	CreateAt   string
	UpdateAt   string
}{
	ID:         "id",
	IDentifier: "identifier",
	ChannelID:  "channel_id",
	AccountID:  "account_id",
	Cost:       "cost",
	Click:      "click",
	Imp:        "imp",
	Ctr:        "ctr",
	Cv:         "cv",
	Cvr:        "cvr",
	Cpa:        "cpa",
	Status:     "status",
	Type:       "type",
	Currency:   "currency",
	RecordDate: "record_date",
	CreateAt:   "create_at",
	UpdateAt:   "update_at",
}
