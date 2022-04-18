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
		headerHandler(req.Headers)
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
		cookies := c.Cookies(url)
		//将新的cookies与旧cookies合并
		mergeCookies(s.Cookies, cookies)
		//更新数据库
		if err := s.UpdateModelCookies(); err != nil {
			e = fmt.Errorf("同步cookie至数据库失败:%w", err)
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
	global.LOG.Info("popin登录成功")
	return
}

func (s *Spider) loginError(err error) error {
	s.Close = true
	panic(fmt.Errorf(errorx.New("popin登录失败", err)))
}
