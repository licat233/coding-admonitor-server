package model

import "time"

const NEWEST = "newest" //最新数据
const NODE = "node"     //节点数据
const DAY = "day"       //全天数据
const LOGIC = "logic"   //逻辑数据

type RecordInter interface {
	ToRecord() *Record
}

type Record struct {
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

//主要用于实现RecordInter接口
func (r *Record) ToRecord() *Record {
	return r
}

//主要用于实现RecordInter接口
func (p *Popin) ToRecord() *Record {
	return &Record{
		ID:         p.ID,
		Account:    p.Account,
		Cost:       p.Cost,
		Click:      p.Click,
		Imp:        p.Imp,
		Ctr:        p.Ctr,
		Cv:         p.Cv,
		Cvr:        p.Cvr,
		Cpa:        p.Cpa,
		Status:     p.Status,
		Type:       p.Type,
		Currency:   p.Currency,
		RecordDate: p.RecordDate,
		CreateAt:   p.CreateAt,
		UpdateAt:   p.UpdateAt,
	}
}

//主要用于实现RecordInter接口
func (l *Line) ToRecord() *Record {
	return &Record{
		ID:         l.ID,
		Account:    l.Account,
		Cost:       l.Cost,
		Click:      l.Click,
		Imp:        l.Imp,
		Ctr:        l.Ctr,
		Cv:         l.Cv,
		Cvr:        l.Cvr,
		Cpa:        l.Cpa,
		Status:     l.Status,
		Type:       l.Type,
		Currency:   l.Currency,
		RecordDate: l.RecordDate,
		CreateAt:   l.CreateAt,
		UpdateAt:   l.UpdateAt,
	}
}
