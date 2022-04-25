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
	baseCookies := []*http.Cookie{
		genHttpCookie("currency", "TWD"),
		genHttpCookie("new_db_trend", "spend"),
		genHttpCookie("db_campaign", "Management#campaignName%2CCost#spend%2CCampaigns#campaignBalance%2CTrafficAcquisition#cpc%2CTrafficAcquisition#impression%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr%2Cconversion#conversion%2Cconversion#conversionRate"),
		genHttpCookie("db_account", "Management#accountName%2CCost#spend%2CTrafficAcquisition#avgCpc%2CTrafficAcquisition#impression%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr%2Cconversion#conversion%2Cconversion#conversionRate%2Cconversion#conversionPrice"),
	}
	s.Cookies = append(s.Cookies, baseCookies...)

	cookies, err := s.CookieModel.FindCookies()
	if err != nil {
		panic(fmt.Errorf("从数据库获取cookie失败:%w", err))
	}
	if len(cookies) == 0 {
		//数据库中不存在cookie，需要登录获取
		s.Login()
		return
	}
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
		Domain:  "ads.popin.cc",
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
		Domain:     "ads.popin.cc",
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

/** mergeCookies 合并cookies到spider.Cookies中，
 * 如果有相同的cookie，则覆盖，
 * 如果没有相同的cookie，则新增到spider.Cookies
 */
func (s *Spider) mergeCookies(newCookies []*http.Cookie) {
	allCookies := append(s.Cookies, newCookies...)
	allM := map[string]string{}
	for _, cookie := range allCookies {
		allM[cookie.Name] = cookie.Value
	}
	cookies := []*http.Cookie{}
	for name, value := range allM {
		cookies = append(cookies, genHttpCookie(name, value))
	}
	s.Cookies = cookies
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
