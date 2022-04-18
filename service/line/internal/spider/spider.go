/*
 * @Author: licat
 * @Date: 2022-04-16 16:13:38
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 18:31:39
 * @Description: loading edit
 */
package spider

import (
	"admonitor-server/service/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type Spider struct {
	Channel         string //渠道名
	Passport        string //账户名
	CookieModel     *model.CookieModel
	CacheStatistics *StatisticsResp //上次相同url请求结果缓存
	CacheUrl        string          //上次相同请求url缓存
	StaticTime      time.Time       //上次爬取statistic时间
	Cookies         []*http.Cookie  //存储cookie，优先从mysql加载
	Domain          string          //请求域
	Close           bool            //状态，注意：只能由login函数进行关闭
	LoginTime       time.Time       //上次登录时间，防止重复登录
	LoginParams     []byte          //账户登录信息
}

func NewSpider() *Spider {
	channel := "line"
	passport := viper.GetString("line.login-info.email")
	line_loginInfo := viper.GetStringMap("line.login-info")
	if len(line_loginInfo) == 0 {
		panic(errors.New("获取line.login-info失败"))
	}
	loginParams, err := json.Marshal(line_loginInfo)
	if err != nil {
		panic(fmt.Errorf("json.Marshal(line_loginInfo) error: %s", err.Error()))
	}

	s := &Spider{
		Channel:         channel,
		Passport:        passport,
		CookieModel:     model.NewCookieModel(channel, passport),
		CacheStatistics: new(StatisticsResp),
		CacheUrl:        "",
		StaticTime:      time.Time{},
		Cookies:         nil,
		Domain:          "account.line.biz",
		Close:           false,
		LoginTime:       time.Time{},
		LoginParams:     loginParams,
	}
	s.InitializeCookie()
	return s
}

func (s *Spider) InitializeCookie() {
	//设置cookies
	s.setCookies()
}
