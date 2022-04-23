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
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

/** accountInfoResp的响应 **/
type accountInfoResp struct {
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
	//通过获取accountInfo检测cookie是否有效
	if s.accountInfo() == nil {
		//cookie失效
		//需要登录，获取新的cookie
		s.Login()
	}
}

/** accountInfo 获取用户信息，用于检测cookie是否有效
 * 如果 spider已经停止，则不执行
 * 如果获取失败，报401错误，则重新Login获取新的cookie
 * 如果获取成功，则更新spider的cookies，并且更新至数据库
 */
func (s *Spider) accountInfo() *accountInfoResp {
	if s.Close {
		return nil
	}
	url := "https://ads.popin.cc/platform/api/campaignManager/accountInfo"
	var e error
	Colly := colly.NewCollector(colly.Async(false))
	Colly.OnRequest(func(req *colly.Request) {
		accountInfoHeader(req.Headers)
	})
	Colly.SetCookies(url, s.Cookies)
	resp := new(accountInfoResp)
	Colly.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = errors.New("响应状态码不是200")
			return
		}
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
		//访问成功
		cookies := Colly.Cookies(url)
		// for _, cookie := range cookies {
		// 	fmt.Printf("%s:%s\n", cookie.Name, cookie.Value)
		// }
		//将新的cookies与旧cookies合并
		s.mergeCookies(cookies)
		//更新数据库
		if err := s.UpdateModelCookies(); err != nil {
			e = fmt.Errorf("同步cookie至数据库失败:%w", err)
		}
	})
	Colly.OnError(func(resp *colly.Response, err error) {
		// cookie := resp.Request.Headers.Get("Cookie")
		// fmt.Println("cookie:", cookie)
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
		if strings.Contains(e.Error(), "Timeout") {
			//超时了，重试
			global.LOG.Warn("popin_accountInfo请求超时，重新发起请求")
			return s.accountInfo()
		}
		//如果是401错误
		if errors.Is(e, AuthError) {
			global.LOG.Error("popin访问accountInfo无权限，请重新登录")
			return nil
		}
		panic(fmt.Errorf(errorx.New("popin获取accountInfo失败", e)))
	}
	return resp
}

func accountInfoHeader(headers *http.Header) {
	headers.Add("authority", "ads.popin.cc")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("accept", "application/json")
	headers.Add("dnt", "1")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("user-agent", UserAgent)
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("referer", "https://ads.popin.cc/platform/new_dashboard")
	headers.Add("accept-language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
}
