/*
 * @Author: licat
 * @Date: 2022-04-16 16:13:38
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 17:22:25
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
	HeartbeatTimer  time.Duration   //心跳时间间隔，设为 5 分钟一次
	Cookies         []*http.Cookie  //存储cookie，优先从redis加载
	Domain          string          //请求域
	Close           bool            //状态，注意：只能由login函数进行关闭
	LoginTime       time.Time       //上次登录时间，防止重复登录
	LoginParams     []byte          //账户登录信息
}

func NewSpider() *Spider {
	channel := "popin"
	passport := viper.GetString("popin.login-info.email")
	popin_loginInfo := viper.GetStringMap("popin.login-info")
	if len(popin_loginInfo) == 0 {
		panic(errors.New("获取popin.login-info失败"))
	}
	loginParams, err := json.Marshal(popin_loginInfo)
	if err != nil {
		panic(fmt.Errorf("json.Marshal(popin_loginInfo) error: %s", err.Error()))
	}
	s := &Spider{
		Channel:         channel,
		Passport:        passport,
		CookieModel:     model.NewCookieModel(channel, passport),
		CacheStatistics: new(StatisticsResp),
		CacheUrl:        "",
		StaticTime:      time.Time{},
		HeartbeatTimer:  5 * time.Minute,
		Cookies:         []*http.Cookie{},
		Domain:          "ads.popin.cc",
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
	//通过心跳检测cookie是否有效
	s.Heartbeat()
}
