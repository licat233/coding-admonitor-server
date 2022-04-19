package model

import (
	"admonitor-server/common/global"
	"errors"

	"gorm.io/gorm"
)

type StandardModel struct {
	*_StandardMgr
}

func NewStandardModel() *StandardModel {
	return &StandardModel{
		_StandardMgr: StandardMgr(global.DB),
	}
}

func (m *StandardModel) InsertOne(standard *Standard) error {
	result := m.DB.WithContext(m.ctx).Model(Standard{}).Omit("id", "create_at", "update_at").Create(standard)
	return result.Error
}

func (m *StandardModel) FindOne(id int64) (*Standard, error) {
	standard := new(Standard)
	result := m.DB.WithContext(m.ctx).Model(Standard{}).First(standard, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return standard, result.Error
}

func (m *StandardModel) FindStandard(channel, account string) (*Standard, error) {
	standard := new(Standard)
	result := m.DB.WithContext(m.ctx).Model(Standard{}).Where("channel = ? AND account = ?", channel, account).First(standard)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return standard, result.Error
}
