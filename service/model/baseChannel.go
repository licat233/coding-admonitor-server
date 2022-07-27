package model

import (
	"admonitor-server/common/global"

	"gorm.io/gorm"
)

type ChannelModel struct {
	*_ChannelMgr
}

func NewChannelModel() *ChannelModel {
	return &ChannelModel{
		_ChannelMgr: ChannelMgr(global.DB),
	}
}

func (m *ChannelModel) baseModel() *gorm.DB {
	return m.DB.WithContext(m.ctx).Model(Channel{})
}

func (m *ChannelModel) FindOne(id int64) (*Channel, error) {
	channel := new(Channel)
	result := m.baseModel().First(channel, id)
	return channel, result.Error
}

func (m *ChannelModel) FindAll() ([]*Channel, error) {
	var channels []*Channel
	result := m.baseModel().Find(&channels)
	return channels, result.Error
}

func (m *ChannelModel) FindByName(name string) (*Channel, error) {
	channel := new(Channel)
	result := m.baseModel().Where("name = ?", name).First(channel)
	return channel, result.Error
}

func (m *ChannelModel) FindAllByStatus(status bool) ([]*Channel, error) {
	var channels []*Channel
	result := m.baseModel().Where("status = ?", status).Find(&channels)
	return channels, result.Error
}

func (m *ChannelModel) InsertOne(channel *Channel) error {
	result := m.baseModel().Omit("id", "create_at", "update_at").Create(channel)
	return result.Error
}

func (m *ChannelModel) InsertAll(channels []*Channel) error {
	result := m.baseModel().Omit("id", "create_at", "update_at").Create(channels)
	return result.Error
}

func (m *ChannelModel) UpdateStatus(ids []int64, status bool) error {
	return m.baseModel().Where("id IN ?", ids).Update("status", status).Error
}
