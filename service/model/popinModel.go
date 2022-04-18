/*
 * @Author: licat
 * @Date: 2022-04-16 14:21:57
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 16:10:25
 * @Description: popin表model
 */
package model

import (
	"admonitor-server/common/global"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PopinModel struct {
	*_PopinMgr
}

func NewPopinModel() *PopinModel {
	return &PopinModel{
		_PopinMgr: PopinMgr(global.DB),
	}
}

func (m *PopinModel) FindOne(id int64) (*Popin, error) {
	line := new(Popin)
	result := m.DB.WithContext(m.ctx).Model(Popin{}).First(line, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return line, result.Error
}

/** FindByDate
 */
func (m *PopinModel) FindByDate(account string, date time.Time) (*Popin, error) {
	line := new(Popin)
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Where("account = ? AND DATE(record_date) = ?", account, date.Local().Format("2006-01-02")).First(line)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return line, result.Error
}

func (m *PopinModel) FindByDateAndType(account string, date time.Time, Type string) (*Popin, error) {
	line := new(Popin)
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Where("account = ? AND DATE(record_date) = ? AND type = ?", account, date.Local().Format("2006-01-02"), Type).First(line)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return line, result.Error
}

func (m *PopinModel) FindsByDateAndType(account string, start, end time.Time, Type string) ([]Popin, error) {
	var lines []Popin
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Where("account = ? AND  DATE(record_date) between ? AND ? AND type = ?", account, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Find(&lines)
	return lines, result.Error
}

func (m *PopinModel) InsertOne(line *Popin) error {
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Omit("id", "create_at", "update_at").Create(line)
	return result.Error
}

func (m *PopinModel) InsertAll(lines []*Popin) error {
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Omit("id", "create_at", "update_at").Create(lines)
	return result.Error
}

// func (m *PopinModel) Save(line *Popin) error {
// 	result := m.DB.WithContext(m.ctx).Model(Popin{}).Omit("id", "create_at", "update_at").Save(line)
// 	return result.Error
// }

func (m *PopinModel) UpdateByDateAndType(line *Popin, date time.Time, Type string) error {
	account := line.Account
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Omit("id", "create_at", "update_at").Where("account = ? AND DATE(record_date) = ? AND type = ?", account, date.Local().Format("2006-01-02"), Type).Save(line)
	return result.Error
}

func (m *PopinModel) Delete(line *Popin) error {
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Delete(line)
	return result.Error
}

func (m *PopinModel) DeleteById(id int64) error {
	line := new(Popin)
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Delete(line, id)
	return result.Error
}

func (m *PopinModel) DeleteByBetweenDate(account string, start, end time.Time, Type string) error {
	line := new(Popin)
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Where("account = ? AND DATE(record_date) between ? AND ? AND type = ?", account, start.Local().Format("2006-01-02"), end.Local().Format("2006-01-02"), Type).Delete(line)
	return result.Error
}

func (m *PopinModel) DeleteAll(ids ...int) error {
	var lines []Popin
	result := m.DB.WithContext(m.ctx).Model(Popin{}).Delete(&lines, ids)
	return result.Error
}
