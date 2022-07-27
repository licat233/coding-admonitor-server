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
	*_CookieMgr
}

func NewCookieModel() *CookieModel {
	return &CookieModel{
		_CookieMgr: CookieMgr(global.DB),
	}
}

func (m *CookieModel) baseModel() *gorm.DB {
	return m.DB.WithContext(m.ctx).Model(Cookie{})
}

func (m *CookieModel) whereModel(channel_id int64) *gorm.DB {
	return m.baseModel().Where("channel_id = ?", channel_id)
}

/** FindOne
 * 通过id查找cookie
 */
func (m *CookieModel) FindOne(id int64) (*Cookie, error) {
	cookie := new(Cookie)
	result := m.baseModel().First(cookie, id)
	return cookie, result.Error
}

/** FindCookies
 * 查找条件下所有的cookie
 * 如果不存在，则返回长度为 0 的切片
 */
func (m *CookieModel) FindCookies(channel_id int64) ([]*Cookie, error) {
	var cookies []*Cookie
	result := m.whereModel(channel_id).Find(&cookies)
	return cookies, result.Error
}

func (m *CookieModel) FindByName(channel_id int64, cookie_name string) (*Cookie, error) {
	var cookie *Cookie
	result := m.baseModel().Where("channel_id = ? AND  name = ?", channel_id, cookie_name).First(&cookie)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return cookie, result.Error
}

func (m *CookieModel) InsertOne(cookie *Cookie) error {
	result := m.baseModel().Omit("id", "create_at", "update_at").Create(cookie)
	return result.Error
}

func (m *CookieModel) InsertAll(cookies []*Cookie) error {
	return m.baseModel().Omit("id", "create_at", "update_at").Create(cookies).Error
}

func (m *CookieModel) UpdateOne(cookie *Cookie) error {
	result := m.baseModel().Where("id =? ", cookie.ID).Omit("id", "create_at", "update_at").Updates(cookie)
	return result.Error
}

func (m *CookieModel) DeleteOne(id int64) error {
	result := m.baseModel().Where("id =? ", id).Delete(&Cookie{})
	return result.Error
}

func (m *CookieModel) DeleteByName(channel_id int64, cookie_name string) error {
	cookie := new(Cookie)
	result := m.baseModel().Where("channel_id = ? AND  name = ?", channel_id, cookie_name).Delete(cookie)
	return result.Error
}

func (m *CookieModel) DeleteAll(ids ...int) error {
	result := m.baseModel().Where("id in ? ", ids).Delete(&[]Cookie{})
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
	res, err := m.FindByName(cookie.ChannelID, cookie.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，则插入
		return m.InsertOne(cookie)
	} else if err != nil { //如果查询出错，则返回错误
		return err
	}
	//存在，则更新value值
	cookie.ID = res.ID
	return m.UpdateOne(cookie)
}
