/*
 * @Author: licat
 * @Date: 2022-04-16 16:21:05
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 17:15:34
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

type LoginResp struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

/** Login 登录popin获取新的cookie
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
	url := "https://ads.popin.cc/platform/api/campaignManager/login"
	var e error
	c := colly.NewCollector()
	c.OnRequest(func(req *colly.Request) {
		loginHeader(req.Headers)
	})
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = fmt.Errorf("login接口响应状态码为%d", r.StatusCode)
			return
		}
		resp := new(LoginResp)
		//反序列化
		err := json.Unmarshal(r.Body, resp)
		if err != nil {
			e = errors.New("json反序列化失败")
			return
		}
		if resp.Code != 0 {
			e = errors.New("账号或者密码错误/丢失")
			return
		}
		//登录成功
		//将登录成功后获取到的cookies传给spider
		cookies := c.Cookies(url)
		//echo
		// fmt.Println("登录获取的cookies")
		// for _, cookie := range cookies {
		// 	fmt.Printf("%s=%s\n", cookie.Name, cookie.Value)
		// }
		// fmt.Println("登录获取的cookies-end")
		//将新的cookies与原cookies合并
		s.mergeCookies(cookies)
		// 此时的cookie还不完整，所以不能存储到数据库中
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
		if strings.Contains(e.Error(), "Timeout") {
			//超时了，重试
			global.LOG.Warn("popin_login请求超时，重新发起请求")
			s.Login()
			return
		}
		s.loginError(e)
		return
	}
	global.LOG.Info("popin登录成功")
	//但仍需要后续的验证，此时获取到的cookie并不完整
	s.authorize()
	//验证成功后，需要验证完整的cookie是否可使用
	if s.accountInfo() == nil {
		s.loginError(errors.New("账号信息获取失败，登录后的cookie不能使用"))
	}
	//更新数据库
	if err := s.UpdateModelCookies(); err != nil {
		s.loginError(fmt.Errorf("同步cookie至数据库失败:%w", err))
	}
	return
}

func (s *Spider) loginError(err error) error {
	s.Close = true
	panic(fmt.Errorf(errorx.New("popin登录失败", err)))
}

func loginHeader(headers *http.Header) {
	headers.Add("authority", "ads.popin.cc")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("accept", "application/json, text/plain, */*")
	headers.Add("dnt", "1")
	headers.Add("content-type", "application/json; charset=UTF-8")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("user-agent", UserAgent)
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("origin", "https://ads.popin.cc")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("referer", "https://ads.popin.cc/platform/user/login")
	headers.Add("accept-language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
}
