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

type StatisticsResp struct {
	Datas  []AccountCostData `json:"datas"`
	Paging struct {
		Page          int64    `json:"page"`
		Size          int64    `json:"size"`
		TotalElements int64    `json:"total_elements"`
		Sorts         []string `json:"sorts"`
	} `json:"paging"`
	ReportCommand struct {
		Since string `json:"since"`
		Until string `json:"until"`
	} `json:"reportCommand"`
}

type AccountCostData struct {
	Id               string `json:"id"`
	AccountInfo      `json:"account"`
	ConfiguredStatus string `json:"configuredStatus"`
	DeliveryStatus   string `json:"deliveryStatus"`
	Statistics       `json:"statistics"`
}

type AccountInfo struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ProductName      string `json:"productName"`
	Timezone         string `json:"timezone"`
	ConfiguredStatus string `json:"configuredStatus"`
	CorporateName    string `json:"corporateName"`
	Country          string `json:"country"`
	Currency         string `json:"currency"`
	CurrencySymbol   string `json:"currencySymbol"`
	DeliveryStatus   string `json:"deliveryStatus"`
}

type Statistics struct {
	Imp            int64   `json:"imp"`
	Click          int64   `json:"click"`
	Cv             int64   `json:"cv"`
	Ctr            float64 `json:"ctr"`
	Cvr            float64 `json:"cvr"`
	Cpm            float64 `json:"cpm"`
	Cpc            float64 `json:"cpc"`
	Cpa            float64 `json:"cpa"`
	Cost           float64 `json:"cost"`
	Currency       string  `json:"currency"`
	CurrencySymbol string  `json:"currencySymbol"`
}

func satisticsURL(date time.Time) string {
	todayDate := date.Format("2006-01-02")
	return fmt.Sprintf("https://admanager.line.biz/api/adservice/v1/my/adaccounts/statistics?groupId=&page=1&size=100&searchKey=&sort=&reportSince=%s&reportUntil=%s", todayDate, todayDate)
}

//Statistics 获取统计数据
//如果 spider已经停止，则返回 nil，
//如果报401错误，则重新Login获取新的cookie
func (s *Spider) Statistics(date time.Time) (resp *StatisticsResp) {
	if s.Close {
		return nil
	}
	url := satisticsURL(date)
	if time.Since(s.StaticTime) < 30*time.Second && s.CacheUrl == url { //短时间内请求，且url相同，不再请求，直接返回缓存数据
		return s.CacheStatistics
	}

	Colly := colly.NewCollector(colly.Async(false))
	Colly.SetRequestTimeout(120 * time.Second)
	Colly.OnRequest(func(req *colly.Request) {
		statisticsHeader(req.Headers)
	})
	Colly.SetCookies(url, s.Cookies)
	var e error
	Colly.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = fmt.Errorf("statistics响应状态码为%d", r.StatusCode)
			return
		}
		resp = new(StatisticsResp)
		//反序列化
		err := json.Unmarshal(r.Body, resp)
		if err != nil {
			e = errors.New("json反序列化失败")
			return
		}
		s.StaticTime = time.Now().Local()
		s.CacheStatistics = resp
		s.CacheUrl = url
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
		if strings.Contains(e.Error(), "Timeout") {
			//超时了，重试
			global.LOG.Warn("line_statistics请求超时，重新发起请求")
			return s.Statistics(date)
		}
		s.statisticsError(e)
		return
	}
	return resp
}

func (s *Spider) statisticsError(err error) {
	//如果是401错误
	if errors.Is(err, AuthError) {
		global.LOG.Error("获取statistic数据无访问权限，准备重新登陆")
		//重新登录，更新cookie
		s.Login()
		return
	}

	panic(fmt.Errorf(errorx.New("line_statistics失败", err)))
}

func statisticsHeader(headers *http.Header) {
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Add("Connection", "keep-alive")
	headers.Add("Host", "admanager.line.biz")
	todayDate := time.Now().Local().Format("2006-01-02")
	referer := fmt.Sprintf("https://admanager.line.biz/adaccount/?duration=hourly&reportSince=%s&reportUntil=%s", todayDate, todayDate)
	headers.Add("Referer", referer)
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "macOS")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Site", "same-origin")
	headers.Add("User-Agent", UserAgent)
	headers.Add("sentry-trace", "3e630a83fc154dcb8603694c5f14dabd-88912279c7b81e26-0")
}
