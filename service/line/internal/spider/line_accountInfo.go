package spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
)

type accountInfoResp struct {
	ID             int64         `json:"id"`
	DisplayID      string        `json:"displayId"`
	UserType       string        `json:"userType"`
	Status         string        `json:"status"`
	BusinessID     string        `json:"businessId"`
	Name           string        `json:"name"`
	Email          string        `json:"email"`
	AdaccountIds   []string      `json:"adaccountIds"`
	GroupIds       []interface{} `json:"groupIds"`
	AdaccountCount int           `json:"adaccountCount"`
	GroupCount     int           `json:"groupCount"`
	Language       string        `json:"language"`
	CreatedDate    time.Time     `json:"createdDate"`
	ModifiedDate   time.Time     `json:"modifiedDate"`
}

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
		//如果是401错误
		if errors.Is(e, AuthError) {
			global.LOG.Error("line访问accountInfo无权限，请重新登录")
			return nil
		}
		panic(fmt.Errorf(errorx.New("line获取accountInfo失败", e)))
	}
	return resp
}

func accountInfoHeader(headers *http.Header) {
	headers.Add("Connection", "keep-alive")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("DNT", "1")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("User-Agent", UserAgent)
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("Sec-Fetch-Site", "same-origin")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://admanager.line.biz/adaccount/")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
}
