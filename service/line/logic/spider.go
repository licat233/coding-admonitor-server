package logic

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
)

type Spider struct {
	CacheData *RespData
	CacheUrl  string
	NodeTime  time.Time
	Timer     time.Duration
	Close     bool
}

func NewSpider() *Spider {
	s := &Spider{
		CacheData: new(RespData),
		CacheUrl:  "",
		NodeTime:  time.Time{},
		Timer:     30 * time.Second,
		Close:     false,
	}
	return s
}

func (s *Spider) Do(date time.Time) (body *RespData) {
	if s.Close {
		return nil
	}
	url := reqURL(date)
	if time.Since(s.NodeTime) < s.Timer && s.CacheUrl == url { //短时间内请求，且url相同
		return s.CacheData
	}

	Colly := colly.NewCollector(colly.Async(false))
	Colly.OnRequest(func(req *colly.Request) {
		headerHandler(req.Headers)
	})

	Colly.OnResponse(func(r *colly.Response) {
		body = new(RespData)
		if err := json.Unmarshal(r.Body, body); err != nil {
			global.LOG.Error(errorx.New("json映射失败", err))
			return
		}
		s.NodeTime = time.Now().Local()
		*s.CacheData = *body
		s.CacheUrl = url
	})
	Colly.OnError(func(resp *colly.Response, err error) {
		if resp.StatusCode == 401 {
			global.LOG.Error(errorx.New("请求失败，无权限访问", err))
			//关闭spider
			s.Close = true
		} else {
			global.LOG.Error(errorx.New("请求失败", err))
		}
	})
	if err := Colly.Visit(url); err != nil {
		global.LOG.Error(errorx.New("Colly爬取数据失败", err))
	}
	Colly.Wait()
	return body
}

func reqURL(date time.Time) string {
	todayDate := date.Format("2006-01-02")
	return fmt.Sprintf("https://admanager.line.biz/api/adservice/v1/my/adaccounts/statistics?groupId=&page=1&size=100&searchKey=&sort=&reportSince=%s&reportUntil=%s", todayDate, todayDate)
}

func refererURL() string {
	todayDate := time.Now().Local().Format("2006-01-02")
	return fmt.Sprintf("https://admanager.line.biz/adaccount/?duration=hourly&reportSince=%s&reportUntil=%s", todayDate, todayDate)
}

func cookieData() string {
	cookie := viper.GetString("line.cookie")
	return cookie
}

func ua() string {
	res := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
	return res
}

func headerHandler(headers *http.Header) {
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Add("Connection", "keep-alive")
	headers.Add("Cookie", cookieData())
	headers.Add("Host", "admanager.line.biz")
	headers.Add("Referer", refererURL())
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "macOS")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Site", "same-origin")
	headers.Add("sentry-trace", "3e630a83fc154dcb8603694c5f14dabd-88912279c7b81e26-0")
}
