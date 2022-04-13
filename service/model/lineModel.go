package model

import (
	"admonitor-server/common/global"
	"time"
)

const NEWEST = "newest" //最新数据
const NODE = "node"     //节点数据
const DAY = "day"       //全天数据
const LOGIC = "logic"   //逻辑数据

type LineModel struct {
	*_LineMgr
}

func NewLineModel() *LineModel {
	return &LineModel{
		_LineMgr: LineMgr(global.DB),
	}
}

func (m *LineModel) FindOne(id int64) (*Line, error) {
	line := new(Line)
	result := m.DB.WithContext(m.ctx).Model(Line{}).First(line, id)
	return line, result.Error
}

func (m *LineModel) FindByDate(account string, date time.Time) (*Line, error) {
	line := new(Line)
	result := m.DB.WithContext(m.ctx).Model(Line{}).Where("account = ? AND DATE(record_date) = ?", account, date.Local().Format("2006-01-02")).First(line)
	return line, result.Error
}

func (m *LineModel) FindByDateAndType(account string, date time.Time, Type string) (*Line, error) {
	line := new(Line)
	result := m.DB.WithContext(m.ctx).Model(Line{}).Where("account = ? AND DATE(record_date) = ? AND type = ?", account, date.Local().Format("2006-01-02"), Type).First(line)
	return line, result.Error
}

func (m *LineModel) FindsByDateAndType(account string, start, end time.Time, Type string) ([]Line, error) {
	var lines []Line
	result := m.DB.WithContext(m.ctx).Model(Line{}).Where("account = ? AND  DATE(record_date) between ? AND ? AND type = ?", account, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Find(&lines)
	return lines, result.Error
}

func (m *LineModel) InsertOne(line *Line) error {
	result := m.DB.WithContext(m.ctx).Model(Line{}).Omit("id", "create_at", "update_at").Create(line)
	return result.Error
}

func (m *LineModel) InsertAll(lines []*Line) error {
	result := m.DB.WithContext(m.ctx).Model(Line{}).Omit("id", "create_at", "update_at").Create(lines)
	return result.Error
}

// func (m *LineModel) Save(line *Line) error {
// 	result := m.DB.WithContext(m.ctx).Model(Line{}).Omit("id", "create_at", "update_at").Save(line)
// 	return result.Error
// }

func (m *LineModel) UpdateByDateAndType(line *Line, date time.Time, Type string) error {
	account := line.Account
	result := m.DB.WithContext(m.ctx).Model(Line{}).Omit("id", "create_at", "update_at").Where("account = ? AND DATE(record_date) = ? AND type = ?", account, date.Local().Format("2006-01-02"), Type).Save(line)
	return result.Error
}

func (m *LineModel) Delete(line *Line) error {
	result := m.DB.WithContext(m.ctx).Model(Line{}).Delete(line)
	return result.Error
}

func (m *LineModel) DeleteById(id int64) error {
	line := new(Line)
	result := m.DB.WithContext(m.ctx).Model(Line{}).Delete(line, id)
	return result.Error
}

func (m *LineModel) DeleteByBetweenDate(account string, start, end time.Time, Type string) error {
	line := new(Line)
	result := m.DB.WithContext(m.ctx).Model(Line{}).Where("account = ? AND DATE(record_date) between ? AND ? AND type = ?", account, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Delete(line)
	return result.Error
}

func (m *LineModel) DeleteAll(ids ...int) error {
	var lines []Line
	result := m.DB.WithContext(m.ctx).Model(Line{}).Delete(&lines, ids)
	return result.Error
}
