package model

import (
	"time"
)

// Cookie 爬虫cookies表
type Cookie struct {
	ID       int64     `gorm:"primaryKey;column:id" json:"-"`    // 主键
	Channel  string    `gorm:"column:channel" json:"channel"`    // cookie归属渠道，如：line、popin
	Account  string    `gorm:"column:account" json:"account"`    // cookie归属账号，如：licat@gmail.com
	Name     string    `gorm:"column:name" json:"name"`          // 设置的cookie_name
	Value    string    `gorm:"column:value" json:"value"`        // 设置的cookie_value
	Domain   string    `gorm:"column:domain" json:"domain"`      // cookie的domain
	Path     string    `gorm:"column:path" json:"path"`          // cookie的path
	CreateAt time.Time `gorm:"column:create_at" json:"createAt"` // 创建时间
	UpdateAt time.Time `gorm:"column:update_at" json:"updateAt"` // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Cookie) TableName() string {
	return "cookie"
}

// CookieColumns get sql column name.获取数据库列名
var CookieColumns = struct {
	ID       string
	Channel  string
	Account  string
	Name     string
	Value    string
	Domain   string
	Path     string
	CreateAt string
	UpdateAt string
}{
	ID:       "id",
	Channel:  "channel",
	Account:  "account",
	Name:     "name",
	Value:    "value",
	Domain:   "domain",
	Path:     "path",
	CreateAt: "create_at",
	UpdateAt: "update_at",
}

// Line Line直投广告数据
type Line struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"-"`        // 主键
	Account    string    `gorm:"column:account" json:"account"`        // 广告账号
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
	RecordDate time.Time `gorm:"column:record_date" json:"recordDate"` // 记录日期
	CreateAt   time.Time `gorm:"column:create_at" json:"createAt"`     // 创建时间
	UpdateAt   time.Time `gorm:"column:update_at" json:"updateAt"`     // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Line) TableName() string {
	return "line"
}

// LineColumns get sql column name.获取数据库列名
var LineColumns = struct {
	ID         string
	Account    string
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
	Account:    "account",
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

// Popin Line直投广告数据
type Popin struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"-"`        // 主键
	Account    string    `gorm:"column:account" json:"account"`        // 广告账号
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
	RecordDate time.Time `gorm:"column:record_date" json:"recordDate"` // 记录日期
	CreateAt   time.Time `gorm:"column:create_at" json:"createAt"`     // 创建时间
	UpdateAt   time.Time `gorm:"column:update_at" json:"updateAt"`     // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Popin) TableName() string {
	return "popin"
}

// PopinColumns get sql column name.获取数据库列名
var PopinColumns = struct {
	ID         string
	Account    string
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
	Account:    "account",
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
