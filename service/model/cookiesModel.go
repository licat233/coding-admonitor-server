/*
 * @Author: licat
 * @Date: 2022-04-16 14:42:01
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 16:54:03
 * @Description: 操作cookie表
 */
package model

import (
	"admonitor-server/common/global"
	"errors"

	"gorm.io/gorm"
)

type CookieModel struct {
	Channel string
	Account string //注意，此处的账号是指登录的账号，通常为邮箱账号，并非子账户
	*_CookieMgr
}

func NewCookieModel(channel, account string) *CookieModel {
	return &CookieModel{
		Channel:    channel,
		Account:    account,
		_CookieMgr: CookieMgr(global.DB),
	}
}

/** FindOne
 * 通过id查找cookie
 */
func (m *CookieModel) FindOne(id int64) (*Cookie, error) {
	cookie := new(Cookie)
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).First(cookie, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return cookie, result.Error
}

/** FindCookies
 * 查找条件下所有的cookie
 * 如果不存在，则返回长度为 0 的切片
 */
func (m *CookieModel) FindCookies() ([]*Cookie, error) {
	var cookies []*Cookie
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("channel = ? AND  account = ?", m.Channel, m.Account).Find(&cookies)
	return cookies, result.Error
}

func (m *CookieModel) FindByName(cookie_name string) (*Cookie, error) {
	var cookie *Cookie
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("channel = ? AND  account = ?", m.Channel, m.Account).Where("name = ?", cookie_name).First(&cookie)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return cookie, result.Error
}

func (m *CookieModel) InsertOne(cookie *Cookie) error {
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Omit("id", "create_at", "update_at").Create(cookie)
	return result.Error
}

func (m *CookieModel) InsertAll(cookies []*Cookie) error {
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Omit("id", "create_at", "update_at").Create(cookies)
	return result.Error
}

func (m *CookieModel) UpdateOne(cookie *Cookie) error {
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("id = ?", cookie.ID).Omit("create_at", "update_at").Save(cookie)
	return result.Error
}

func (m *CookieModel) Delete(cookie *Cookie) error {
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("channel = ? AND  account = ?", m.Channel, m.Account).Delete(cookie)
	return result.Error
}

func (m *CookieModel) DeleteByName(cookie_name string) error {
	cookie := new(Cookie)
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("channel = ? AND  account = ?", m.Channel, m.Account).Where("name = ?", cookie_name).Delete(cookie)
	return result.Error
}

func (m *CookieModel) DeleteById(id int64) error {
	cookie := new(Cookie)
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("channel = ? AND  account = ?", m.Channel, m.Account).Delete(cookie, id)
	return result.Error
}

func (m *CookieModel) DeleteAll(ids ...int) error {
	var cookies []Cookie
	result := m.DB.WithContext(m.ctx).Model(Cookie{}).Where("channel = ? AND  account = ?", m.Channel, m.Account).Delete(&cookies, ids)
	return result.Error
}

//Replaces 批量替换，不存在则创建
func (m *CookieModel) Replaces(cookies []*Cookie) error {
	for _, cookie := range cookies {
		if err := m.Replace(cookie); err != nil {
			return err
		}
	}
	return nil
}

//Replace 替换，不存在则创建
func (m *CookieModel) Replace(cookie *Cookie) error {
	res, err := m.FindByName(cookie.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，则插入
		if err := m.InsertOne(cookie); err != nil {
			//插入失败
			return err
		}
		//插入成功
		return nil
	}
	//如果查询出错，则返回错误
	if err != nil {
		return err
	}
	//存在，则更新value值
	res.Value = cookie.Value
	if err := m.UpdateOne(res); err != nil {
		return err
	}
	return nil
}
