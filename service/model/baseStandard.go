package model

import (
	"admonitor-server/common/global"

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

func (m *StandardModel) baseModel() *gorm.DB {
	return m.DB.WithContext(m.ctx).Model(Standard{})
}

func (m *StandardModel) InsertOne(standard *Standard) error {
	result := m.baseModel().Omit("id", "create_at", "update_at").Create(standard)
	return result.Error
}

func (m *StandardModel) FindOne(id int64) (*Standard, error) {
	standard := new(Standard)
	result := m.baseModel().First(standard, id)
	return standard, result.Error
}

func (m *StandardModel) FindByAccountId(account_id int64) (*Standard, error) {
	standard := new(Standard)
	result := m.baseModel().Where("account_id = ?", account_id).First(standard)
	return standard, result.Error
}
