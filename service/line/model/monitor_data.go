package model

import (
	"time"
)

// Line Line直投广告数据
type Line struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"-"`        // 主键
	Account    string    `gorm:"column:account" json:"account"`        // 广告账号
	Cost       float64   `gorm:"column:cost" json:"cost"`              // 消耗额
	Click      uint      `gorm:"column:click" json:"click"`            // 点击数
	Imp        uint      `gorm:"column:imp" json:"imp"`                // 展现量
	Ctr        float64   `gorm:"column:ctr" json:"ctr"`                // 点击率
	Cv         uint      `gorm:"column:cv" json:"cv"`                  // 转化数
	Cvr        float64   `gorm:"column:cvr" json:"cvr"`                // 转化率
	Cpa        float64   `gorm:"column:cpa" json:"cpa"`                // 转化成本
	Status     uint8     `gorm:"column:status" json:"status"`          // 状态
	Type       uint8     `gorm:"column:type" json:"type"`              // 数据类型：0表示最新数据，1表示全天数据，2表示时间节点数据
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
