package internal

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
	CacheData      *RespData
	CacheUrl       string
	NodeTime       time.Time
	Timer          time.Duration
	HeartbeatTimer time.Duration
	Close          bool
}

func NewSpider() *Spider {
	s := &Spider{
		CacheData:      new(RespData),
		CacheUrl:       "",
		NodeTime:       time.Time{},
		Timer:          30 * time.Second,
		HeartbeatTimer: 5 * time.Minute,
		Close:          false,
	}
	return s
}

//Do 开始爬取数据
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
	dateStr := date.Format("20060102")
	return fmt.Sprintf("https://ads.popin.cc/platform/api/campaignManager/dashboard/getAccountDashboardList?dateType=1&timezone=UTC%2B8&startDate=%s&endDate=%s&currency=USD&accountName=&accountIds=&currentPage=1&pageSize=50&sortField=&sortVal=", dateStr, dateStr)
}

func cookieData() string {
	cookie := viper.GetString("popin.cookie")
	return cookie
}

func ua() string {
	res := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
	return res
}

func headerHandler(headers *http.Header) {
	headers.Add("Accept", "application/json")
	headers.Add("User-Agent", ua())
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
	headers.Add("Connection", "keep-alive")
	headers.Add("Cookie", cookieData())
	headers.Add("Referer", "https://ads.popin.cc/platform/new_dashboard")
	headers.Add("Referrer-Policy", "strict-origin-when-cross-origin")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "macOS")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Site", "same-origin")
}

//Heartbeat 心跳
func (s *Spider) Heartbeat() {
	if s.Close {
		return
	}
	var e error
	Colly := colly.NewCollector(colly.Async(false))
	Colly.OnRequest(func(req *colly.Request) {
		headerHandler(req.Headers)
		req.Headers.Set("Cookie", heartbeatCookie())
	})

	Colly.OnResponse(func(r *colly.Response) {
		// body := new(HeartbeatResp)
		// if err := json.Unmarshal(r.Body, body); err != nil {
		// 	global.LOG.Error(errorx.New("json映射失败", err))
		// 	e = fmt.Errorf("json映射失败%w", err)
		// 	return
		// }
	})
	Colly.OnError(func(resp *colly.Response, err error) {
		if resp.StatusCode == 401 {
			global.LOG.Error(errorx.New("请求失败，无权限访问", err))
			e = fmt.Errorf("请求失败，无权限访问%w", err)
			//关闭spider
			s.Close = true
		} else {
			global.LOG.Error(errorx.New("请求失败", err))
			e = fmt.Errorf("请求失败%w", err)
		}
	})
	if err := Colly.Visit("https://ads.popin.cc/platform/api/campaignManager/accountInfo"); err != nil {
		global.LOG.Error(errorx.New("Colly请求失败", err))
		e = fmt.Errorf("Colly请求失败%w", err)
	}
	Colly.Wait()
	if e != nil {
		s.Close = true
		global.LOG.Error(errorx.New("popin心跳错误，已关闭spider", e))
	}
	return
}

func heartbeatCookie() string {
	return viper.GetString("popin.heartbeat_cookie")
}
