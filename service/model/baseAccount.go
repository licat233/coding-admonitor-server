package model

import (
	"admonitor-server/common/global"
	"errors"

	"gorm.io/gorm"
)

type AccountModel struct {
	*_AccountMgr
}

func NewAccountModel() *AccountModel {
	return &AccountModel{
		_AccountMgr: AccountMgr(global.DB),
	}
}

func (m *AccountModel) baseModel() *gorm.DB {
	return m.DB.WithContext(m.ctx).Model(Account{})
}

func (m *AccountModel) FindOne(id int64) (*Account, error) {
	account := new(Account)
	result := m.baseModel().First(account, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return account, result.Error
}

func (m *AccountModel) FindAll() ([]*Account, error) {
	var accounts []*Account
	result := m.baseModel().Find(&accounts)
	return accounts, result.Error
}

func (m *AccountModel) FindByAid(aid string) (*Account, error) {
	var account = new(Account)
	result := m.baseModel().Where("aid = ?", aid).First(account)
	return account, result.Error
}

func (m *AccountModel) FindAllByStatus(status bool) ([]*Account, error) {
	var accounts []*Account
	result := m.baseModel().Where("status = ?", status).Find(&accounts)
	return accounts, result.Error
}

func (m *AccountModel) FindAllByChannelId(channel_id int64) ([]*Account, error) {
	var accounts []*Account
	result := m.baseModel().Where("channel_id = ?", channel_id).Find(&accounts)
	return accounts, result.Error
}

func (m *AccountModel) FindAllByChannelAndStatus(channel_id int64, status bool) (*Account, error) {
	var accounts *Account
	result := m.baseModel().Where("status = ? AND channel_id = ?", status, channel_id).Find(&accounts)
	return accounts, result.Error
}

func (m *AccountModel) FindAllByChannelsAndStatus(channel_ids []int64, status bool) ([]*Account, error) {
	var accounts []*Account
	result := m.baseModel().Where("status = ? AND channel_id IN ?", status, channel_ids).Find(&accounts)
	return accounts, result.Error
}

func (m *AccountModel) InsertOne(account *Account) error {
	return m.baseModel().Omit("id", "create_at", "update_at").Create(account).Error
}

func (m *AccountModel) InsertAll(accounts []*Account) error {
	result := m.baseModel().Omit("id", "create_at", "update_at").Create(accounts)
	return result.Error
}

func (m *AccountModel) UpdateOne(account *Account) error {
	result := m.baseModel().Where("id = ?", account.ID).Omit("id", "create_at", "update_at").Updates(account)
	return result.Error
}

func (m *AccountModel) UpdateStatus(ids []int64, status bool) error {
	return m.baseModel().Where("id IN ?", ids).Update("status", status).Error
}

//查询是否存在
func (m *AccountModel) IsExist(aid string) (bool, error) {
	var count int64
	err := m.baseModel().Where("aid = ?", aid).Limit(1).Count(&count).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

//插入，通过url查询是否存在，如果存在，则不再插入
func (m *AccountModel) InsertByAid(account *Account) error {
	exist, err := m.IsExist(account.Aid)
	if err != nil {
		return err
	}
	if exist {
		//存在，不存储
		return nil
	}
	return m.InsertOne(account)
}

//插入，通过url查询是否存在，如果存在，则更新
func (m *AccountModel) Replace(account *Account) error {
	exist, err := m.IsExist(account.Aid)
	if err != nil {
		return err
	}
	if exist {
		//存在，更新
		return m.UpdateOne(account)
	}
	return m.InsertOne(account)
}

//ReplaceAll 替换所有，存在则更新，不存在则插入
func (m *AccountModel) ReplaceAll(accounts []*Account) (err error) {
	if len(accounts) == 0 {
		return
	}

	for _, account := range accounts {
		if err = m.Replace(account); err != nil {
			return
		}
	}
	return
}
