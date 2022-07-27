package spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/service/model"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

//InitSpiderCookies 初始化Spider的Cookies
func (s *Spider) InitSpiderCookies() {
	httpCookies := s.GetCookiesFromDB()
	s.MergeCookies(httpCookies)
}

//GenModelCookie 生成model.Cookie
func (s *Spider) GenModelCookie(cookie *http.Cookie) *model.Cookie {
	return &model.Cookie{
		ChannelID: s.Channel.ID,
		Name:      cookie.Name,
		Value:     cookie.Value,
		Domain:    s.Channel.Domain,
		Path:      "/", //统一为根目录，最高权限
	}
}

//GenHttpCookie 生成http.Cookie
func (s *Spider) GenHttpCookie(name, value string) *http.Cookie {
	return GenHttpCookie(name, value, s.Channel.Domain)
}

//GenHttpCookie 生成http.Cookie
func GenHttpCookie(name, value, domain string) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Domain:  domain,
		Expires: time.Now().AddDate(1, 0, 0),
		MaxAge:  31536000,
	}
}

/** UpdateModelCookies
 * 更新spider携带的cookie到数据库
 * 如果已存在，则更新，不存在，不存在则创建
 */
func (s *Spider) UpdateModelCookies() error {
	modelCookies := []*model.Cookie{}
	cookies := s.Cookies
	for _, cookie := range cookies {
		modelCookies = append(modelCookies, s.GenModelCookie(cookie))
	}
	err := s.CookieModel.Replaces(modelCookies)
	if err != nil {
		global.LOG.Error(errorx.New("mysql更新cookies数据失败", err))
	}
	return err
}

/** MergeCookies 合并cookies到spider.Cookies中，
 * 如果有相同的cookie，则覆盖，
 * 如果没有相同的cookie，则新增到spider.Cookies
 */
func (s *Spider) MergeCookies(newCookies []*http.Cookie) {
	allCookies := append(s.Cookies, newCookies...)
	allM := map[string]*http.Cookie{}
	for _, cookie := range allCookies {
		allM[cookie.Name] = cookie
	}
	cookies := []*http.Cookie{}
	for _, cookie := range allM {
		cookies = append(cookies, cookie)
	}
	s.Cookies = cookies
}

func (s *Spider) ContainsCookie(newCookie *http.Cookie) bool {
	for _, cookie := range s.Cookies {
		if cookie.Name == newCookie.Name && cookie.Value == newCookie.Value {
			return true
		}
	}
	return false
}

func (s *Spider) ContainsCookies(newCookies []*http.Cookie) bool {
	for _, cookie := range newCookies {
		if s.ContainsCookie(cookie) {
			return true
		}
	}
	return false
}

//GetCookiesFromDB 从数据库中获取[]*http.cookie
func (s *Spider) GetCookiesFromDB() []*http.Cookie {
	cookies, err := s.CookieModel.FindCookies(s.Channel.ID)
	if len(cookies) == 0 || errors.Is(err, gorm.ErrRecordNotFound) {
		//数据库中不存在cookie
		return nil
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询cookie数据失败", err))
		panic(fmt.Errorf("mysql查询cookie数据失败:%w", err))
	}
	httpCookies := []*http.Cookie{}
	for _, cookie := range cookies {
		httpCookies = append(httpCookies, GenHttpCookie(cookie.Name, cookie.Value, cookie.Domain))
	}
	return httpCookies
}

//GetCookiesFromSpider 从spider中获取[]*http.cookie
func (s *Spider) GetCookiesFromSpider() []*model.Cookie {
	cookies := s.Cookies
	modelCookies := []*model.Cookie{}
	for _, cookie := range cookies {
		modelCookies = append(modelCookies, s.GenModelCookie(cookie))
	}
	return modelCookies
}

//验证spider的cookies是否符合要求
func (s *Spider) ValidateCookies(cookies []*http.Cookie) error {
	if len(cookies) == 0 {
		return errors.New("spider的cookies为空")
	}
	return nil
}

//ModelCookiesToString 多条demolCookies转字符串
func ModelCookiesToString(cookies []*model.Cookie) string {
	if cookies == nil {
		return ""
	}
	var cookie_arr []string
	for _, c := range cookies {
		cookie_arr = append(cookie_arr, c.Name+"="+c.Value)
	}
	cookie_str := strings.Join(cookie_arr, ";")
	return cookie_str
}

//HttpCookiesToString 多条httpCookies转字符串
func HttpCookiesToString(cookies []*http.Cookie) string {
	if len(cookies) == 0 {
		return ""
	}
	var cookie_arr []string
	for _, c := range cookies {
		cookie_arr = append(cookie_arr, c.Name+"="+c.Value)
	}
	cookie_str := strings.Join(cookie_arr, ";")
	return cookie_str
}

//HttpCookiesToString 多条httpCookies转字符串
func (s *Spider) HttpCookiesToString() string {
	return HttpCookiesToString(s.Cookies)
}
