/*
 * @Author: licat
 * @Date: 2022-04-16 14:59:52
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 17:23:45
 * @Description: cookie操作
 */

package spider

import (
	"admonitor-server/service/model"
	"fmt"
	"net/http"
	"time"
)

/** setCookies
 * 数据库加载cookie到spider
 */
func (s *Spider) setCookies() {
	cookies, err := s.CookieModel.FindCookies()
	if err != nil {
		panic(fmt.Errorf("从数据库获取cookie失败:%w", err))
	}
	s.Cookies = []*http.Cookie{}
	for _, cookie := range cookies {
		hcookie := genHttpCookie(cookie.Name, cookie.Value)
		s.Cookies = append(s.Cookies, hcookie)
	}
}

func (s *Spider) genModelCookie(name, value string) *model.Cookie {
	return &model.Cookie{
		ID:      0,
		Channel: s.Channel,
		Account: s.Passport,
		Name:    name,
		Value:   value,
		Domain:  "line.biz",
		Path:    "/",
	}
}

/** UpdateModelCookies
 * 更新spider携带的cookie到数据库
 * 如果已存在，则更新，不存在，不存在则创建
 */
func (s *Spider) UpdateModelCookies() error {
	cookies := s.Cookies
	for _, cookie := range cookies {
		modelCookie := s.genModelCookie(cookie.Name, cookie.Value)
		if err := s.CookieModel.Replace(modelCookie); err != nil {
			return err
		}
	}
	return nil
}

func genHttpCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:       name,
		Value:      value,
		Path:       "/",
		Domain:     "line.biz",
		Expires:    time.Now().AddDate(10, 0, 0),
		RawExpires: "",
		MaxAge:     31536000,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   []string{},
	}
}

var AuthError = fmt.Errorf("请求401错误，无权限访问，cookies失效")

/** mergeCookies 合并两个cookie列表，
 * 如果有相同的cookie，则覆盖，
 * 如果没有相同的cookie，则新增到oldcookies中去
 */
func mergeCookies(oldCookies []*http.Cookie, newCookies []*http.Cookie) {
	for _, newCookie := range newCookies {
		isExist := false
		for _, oldCookie := range oldCookies {
			if newCookie.Name == oldCookie.Name {
				isExist = true
				oldCookie.Value = newCookie.Value
				break
			}
		}
		if isExist != true {
			oldCookies = append(oldCookies, newCookie)
		}
	}
}

func cookiesToString(cookies []*model.Cookie) string {
	if cookies == nil {
		return ""
	}
	var result string
	for _, cookie := range cookies {
		result += fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value)
	}
	return result
}
