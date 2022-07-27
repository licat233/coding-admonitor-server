package popin_spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func (s *PopinSpider) GetCookies() {
	if !s.AllowLogin() {
		return
	}
	reqUrl := "https://ads.popin.cc/platform/api/campaignManager/login"
	var e error
	c := colly.NewCollector()
	c.SetCookies(reqUrl, s.Cookies)
	//开始登录请求
	c.OnRequest(func(req *colly.Request) {
		s.addHeader(req.Headers)
	})
	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return
		}
		resp := new(LoginResp)
		//反序列化
		err := json.Unmarshal(r.Body, resp)
		if err != nil {
			e = errors.New("json反序列化失败")
			global.LOG.Error(e)
			return
		}
		if resp.Code != 0 {
			e = errors.New("账号或者密码错误/丢失")
			return
		}
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = s.OnErrorHandler(resp, err)
	})
	c.PostRaw(reqUrl, s.getLoginRequestData())
	if e != nil {
		s.SpiderErrorHandler(e, s.GetCookies, func() {})
		return
	}
	c.Wait()
	// global.LOG.Info("popin登录成功")
	//开始鉴权请求
	c.OnRequest(func(r *colly.Request) {
		s.addHeader(r.Headers)
	})
	c.OnResponse(func(r *colly.Response) {})
	c.Visit("https://ads.popin.cc/platform/main/login?type=m&source=us") //这里的返回的error与onError里面的error一致
	if e != nil {
		s.SpiderErrorHandler(e, s.GetCookies, s.GetCookies)
	}
	// global.LOG.Info("popin鉴权成功")
	//开始验证请求
	c.OnRequest(func(req *colly.Request) {
		s.accountInfoHeader(req.Headers)
	})
	c.Visit("https://ads.popin.cc/platform/api/campaignManager/accountInfo") //这里的返回的error与onError里面的error一致
	if e != nil {
		s.SpiderErrorHandler(e, s.GetCookies, s.GetCookies)
	}
	// global.LOG.Info("popin验证成功")

	//保存cookies到spider中
	s.Cookies = c.Cookies(reqUrl)
	//开始保存cookies到数据库，存在则更新
	if e = s.UpdateModelCookies(); e != nil {
		msg := errorx.New("同步popin_cookie至数据库失败", e)
		global.LOG.Error(msg)
		panic(errors.New(msg))
	}
	global.LOG.Info("popin登录成功")
}

func (s *PopinSpider) getLoginRequestData() []byte {
	email := s.Channel.Email
	password := s.Channel.Password
	if len(email) == 0 || len(password) == 0 {
		msg := fmt.Sprintf("【%s】渠道账号登录信息未完善", s.Channel.Name)
		global.LOG.Warn(msg)
		panic(errors.New(msg))
	}
	popin_loginInfo := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	loginParams, err := json.Marshal(popin_loginInfo)
	if err != nil {
		panic(fmt.Errorf("json.Marshal(popin_loginInfo) error: %s", err.Error()))
	}
	return loginParams
}
