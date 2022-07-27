package line_spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func (s *LineSpider) GetCookies() {
	if !s.AllowLogin() {
		return
	}
	c := colly.NewCollector()
	var e error
	//访问登录页面
	pageUrl := "https://account.line.biz/login?redirectUri=https%253A%252F%252Fadmanager.line.biz%252Flogin%252F&scope=email"
	c.SetCookies(pageUrl, s.Cookies)
	c.OnRequest(func(req *colly.Request) {
		loginPageHeader(req.Headers)
	})
	var csrfToken_name string
	var csrfToken_value string
	c.OnHTML("head", func(e *colly.HTMLElement) {
		csrfToken_name = e.ChildAttr("meta[name=x-csrf-header]", "content")
		csrfToken_value = e.ChildAttr("meta[name=x-csrf]", "content")
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = s.OnErrorHandler(resp, err)
	})
	c.Visit(pageUrl) //这里的返回的error与onError里面的error一致e
	if e != nil {
		s.SpiderErrorHandler(e, s.GetCookies, func() {})
		return
	}
	//访问登录接口
	loginUrl := "https://account.line.biz/api/login/email"
	c.OnRequest(func(req *colly.Request) {
		loginHeader(req.Headers)
		req.Headers.Add(csrfToken_name, csrfToken_value)
	})

	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return
		}
		statusResp := new(LoginResp)
		//反序列化
		err := json.Unmarshal(r.Body, statusResp)
		if err != nil {
			e = errors.New("json反序列化失败")
			return
		}
		status := statusResp.Status
		if status == "needCaptchaVerification" { //需要验证码
			//TODO:验证码处理，暂未能破解
			//将error定义为timeout，以触发重新登录
			e = errors.New("需要验证码处理，但暂未能破解，定义为 Timeout ，以触发重新登录")
			return
		}
		if status != "success" {
			e = errors.New("账号密码错误")
			return
		}
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = s.OnErrorHandler(resp, err)
	})
	c.PostRaw(loginUrl, s.loginRequestData()) //这里的返回的error与onError里面的error一致
	if e != nil {
		s.SpiderErrorHandler(e, s.GetCookies, func() {})
		return
	}

	//保存cookies到spider中
	s.Cookies = c.Cookies(pageUrl)
	//开始保存cookies到数据库，存在则更新
	if e = s.UpdateModelCookies(); e != nil {
		msg := errorx.New("同步line_cookie至数据库失败", e)
		global.LOG.Error(msg)
		panic(errors.New(msg))
	}
	global.LOG.Info("line登录成功")
}

func (s *LineSpider) loginRequestData() []byte {
	email := s.Channel.Email
	password := s.Channel.Password
	if len(email) == 0 || len(password) == 0 {
		msg := fmt.Sprintf("【%s】渠道账号登录信息未完善", s.Channel.Name)
		global.LOG.Warn(msg)
		panic(errors.New(msg))
	}
	line_loginInfo := map[string]interface{}{
		"email":        email,
		"password":     password,
		"stayLoggedIn": "true",
	}
	loginParams, err := json.Marshal(line_loginInfo)
	if err != nil {
		s.Close = true
		panic(fmt.Errorf("json.Marshal(line_loginInfo) error: %s", err.Error()))
	}
	return loginParams
}
