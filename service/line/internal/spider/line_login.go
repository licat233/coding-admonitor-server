/*
 * @Author: licat
 * @Date: 2022-04-16 16:21:05
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 18:39:15
 * @Description: loading edit
 */
package spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type statusResp struct {
	Status string `json:"status"`
}

//statusCode = 400
type loginFaild struct {
	Status       string         `json:"status"` //failed
	Notices      []interface{}  `json:"notices"`
	GlobalErrors []GlobalErrors `json:"globalErrors"`
	FieldErrors  []interface{}  `json:"fieldErrors"`
}

type GlobalErrors struct {
	Message string `json:"message"`
}

//stutusCode = 200，需要验证码
type captachResp struct {
	Status          string         `json:"status"` //needCaptchaVerification
	CaptchaKey      string         `json:"captchaKey"`
	CaptchaImageURI string         `json:"captchaImageUri"`
	GlobalErrors    []GlobalErrors `json:"globalErrors"`
}

//statusCode = 200
type loginSuccess struct {
	Status      string `json:"status"` //success
	RedirectURI string `json:"redirectUri"`
}

/** Login 登录line获取新的cookie
 * 如果spider已经停止，则不执行
 * 如果登录失败，则panic，上层恢复函数有log处理，这里无需打LOG
 * 如果登录成功，则更新spider的cookies，并且更新至数据库
 */
func (s *Spider) Login() {
	//30秒内不能重复登录
	if time.Since(s.LoginTime) < 30*time.Second {
		global.LOG.Warn("请勿在30秒内重复登录")
		return
	}
	s.LoginTime = time.Now()
	url := "https://account.line.biz/api/login/email"
	//step 1: 访问登录页面，获取Token
	pageCookies, csrfToken_name, csrfToken_value := s.loginPage()
	//step 2: 访问登录接口
	c := colly.NewCollector(colly.Async(false))
	var e error
	//添加请求头信息
	c.OnRequest(func(req *colly.Request) {
		loginHeader(req.Headers)
		//X-XSRF-TOKEN: f1a81ead-e82f-47c5-b091-04159061b6ce
		req.Headers.Add(csrfToken_name, csrfToken_value)
	})
	//设置cookies
	if err := c.SetCookies(url, pageCookies); err != nil {
		s.loginError(err)
	}
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = respError(r)
			return
		}
		statusResp := new(statusResp)
		//反序列化
		err := json.Unmarshal(r.Body, statusResp)
		if err != nil {
			e = errors.New("json反序列化失败")
			return
		}
		status := statusResp.Status
		if status == "needCaptchaVerification" { //需要验证码
			//TODO:验证码处理
			// fmt.Println(r.Headers)
			e = respError(r)
			return
		}
		if status == "success" {
			global.LOG.Info("line登录成功")
			//登录成功，获取该网站的cookies
			// fmt.Println(r.Headers.Get("Set-Cookie"))
			// cookies := c.Cookies(url)
			cookies := []*http.Cookie{}
			//注意:line的cookie故意以==结尾，所以colly自动提取到的cookie值结尾会少一个=号，保险起见，我们自己提取
			cookieValues := r.Headers.Values("Set-Cookie")
			for _, value := range cookieValues {
				//查找第一个【;】位置
				index1 := strings.Index(value, ";")
				if index1 == -1 {
					continue
				}
				//截取字符串
				s1 := value[:index1]
				//查找第一个【=】位置
				index2 := strings.Index(s1, "=")
				if index2 == -1 {
					continue
				}
				//截取字符串
				name := s1[:index2]
				value := s1[index2+1:]
				cookies = append(cookies, genHttpCookie(name, value))
			}
			//更新spider的cookies
			s.mergeCookies(cookies)
			//更新数据库
			if err := s.UpdateModelCookies(); err != nil {
				e = fmt.Errorf("同步cookie至数据库失败:%w", err)
			}
		}
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = fmt.Errorf("访问login接口失败【状态码:%d】【内容:%s】:%w", resp.StatusCode, string(resp.Body), err)
	})
	if err := c.PostRaw(url, s.LoginParams); err != nil {
		s.loginError(err)
		return
	}
	c.Wait()
	if e != nil {
		s.loginError(e)
		return
	}
	global.LOG.Info("line登录成功")
	//即使登录成功，也要使用accountInfo()检测一下新的cookie是否可用
	if s.accountInfo() == nil {
		//如果不能使用，则停止爬虫
		s.loginError(errors.New("新cookie无效"))
	}
	return
}

func respError(r *colly.Response) error {
	return fmt.Errorf("【状态码:%d】【内容:%s】", r.StatusCode, string(r.Body))
}

func (s *Spider) loginError(err error) error {
	s.Close = true
	panic(fmt.Errorf(errorx.New("line登录失败", err)))
}

/** loginPage
 * 访问登录页面
 * 目的：获取TOKEN
 */
func (s *Spider) loginPage() (cookies []*http.Cookie, csrfToken_name, csrfToken_value string) {
	url := "https://account.line.biz/login?redirectUri=https%253A%252F%252Fadmanager.line.biz%252Flogin%252F&scope=email"
	c := colly.NewCollector(colly.Async(false))
	c.OnRequest(func(req *colly.Request) {
		loginPageHeader(req.Headers)
	})
	var e error
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = fmt.Errorf("loginPage响应状态码为%d", r.StatusCode)
			return
		}
		//获取响应的cookie
		cookies = c.Cookies(url)
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = err
	})
	c.OnHTML("head", func(e *colly.HTMLElement) {
		csrfToken_name = e.ChildAttr("meta[name=x-csrf-header]", "content")
		csrfToken_value = e.ChildAttr("meta[name=x-csrf]", "content")
	})
	if err := c.Visit(url); err != nil {
		e = fmt.Errorf("访问loginPage失败:%w", err)
	}
	if e != nil {
		s.loginError(e)
	}
	c.Wait()
	if csrfToken_name == "" || csrfToken_value == "" || len(cookies) == 0 {
		s.loginError(errors.New("访问login页面获取TOKEN失败"))
	}
	return
}

func loginPageHeader(headers *http.Header) {
	headers.Add("Connection", "keep-alive")
	headers.Add("Cache-Control", "max-age=0")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("DNT", "1")
	headers.Add("Upgrade-Insecure-Requests", "1")
	headers.Add("User-Agent", UserAgent)
	headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	headers.Add("Sec-Fetch-Site", "none")
	headers.Add("Sec-Fetch-Mode", "navigate")
	headers.Add("Sec-Fetch-User", "?1")
	headers.Add("Sec-Fetch-Dest", "document")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
}

func loginHeader(headers *http.Header) {
	headers.Add("Connection", "keep-alive")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("DNT", "1")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("User-Agent", UserAgent)
	headers.Add("Content-Type", "application/json;charset=UTF-8")
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("Origin", "https://account.line.biz")
	headers.Add("Sec-Fetch-Site", "same-origin")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://account.line.biz/login?redirectUri=https%3A%2F%2Fadmanager.line.biz%2Flogin%2F&scope=email")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
}
