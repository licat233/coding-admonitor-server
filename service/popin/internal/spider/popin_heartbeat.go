/*
 * @Author: licat
 * @Date: 2022-04-16 16:15:16
 * @LastEditors: licat
 * @LastEditTime: 2022-04-16 17:45:32
 * @Description: loading edit
 */
package spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gocolly/colly/v2"
)

/** 心跳的响应 **/
type HeartbeatResp struct {
	Code int `json:"code"`
	Data struct {
		AdAuditNum             string        `json:"adAuditNum"`
		Balance                string        `json:"balance"`
		CreditLineAccountAlarm string        `json:"creditLineAccountAlarm"`
		SspInfo                []interface{} `json:"sspInfo"`
		UserTags               []UserTags    `json:"userTags"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type UserTags struct {
	Label string `json:"label"`
	Uv    int    `json:"uv"`
	Value string `json:"value"`
}

/* Heartbeat 通过心跳检测cookie是否有效
 * 如果失效，则重新登录，更新cookie
 */
func (s *Spider) Heartbeat() {
	//通过心跳检测cookie是否有效
	if s.heartbeat() == false {
		//cookie失效
		//需要登录，获取新的cookie
		s.Login()
	}
}

/** heartbeat 心跳，用于检测cookie是否有效
 * 如果 spider已经停止，则不执行
 * 如果心跳失败，报401错误，则重新Login获取新的cookie
 * 如果心跳成功，则更新spider的cookies，并且更新至数据库
 */
func (s *Spider) heartbeat() bool {
	if s.Close {
		return false
	}
	url := "https://ads.popin.cc/platform/api/campaignManager/accountInfo"
	var e error
	Colly := colly.NewCollector(colly.Async(false))
	Colly.OnRequest(func(req *colly.Request) {
		headerHandler(req.Headers)
	})
	Colly.SetCookies(s.Domain, s.Cookies)
	Colly.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = errors.New("响应状态码不是200")
			return
		}
		resp := new(HeartbeatResp)
		//反序列化
		err := json.Unmarshal(r.Body, resp)
		if err != nil {
			e = errors.New("json反序列化失败")
			return
		}
		if resp.Code != 0 {
			e = errors.New(string(r.Body))
			return
		}
		//心跳成功
		cookies := Colly.Cookies(url)
		//将新的cookies与旧cookies合并
		mergeCookies(s.Cookies, cookies)
		//更新数据库
		if err := s.UpdateModelCookies(); err != nil {
			e = fmt.Errorf("同步cookie至数据库失败:%w", err)
		}
	})
	Colly.OnError(func(resp *colly.Response, err error) {
		if resp.StatusCode == 401 {
			e = AuthError
		} else {
			e = err
		}
	})
	if err := Colly.Visit(url); err != nil {
		e = err
	}
	Colly.Wait()
	if e != nil {
		//如果是401错误
		if errors.Is(e, AuthError) {
			global.LOG.Error("popin心跳失败，重新登录")
			return false
		}
		panic(fmt.Errorf(errorx.New("popin_Hearbeat失败", e)))
	}
	return true
}
