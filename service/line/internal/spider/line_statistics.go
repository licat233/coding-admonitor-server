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
	Colly.OnRequest(func(req *colly.Request) {
		statisticsHeader(req.Headers)
	})
	Colly.SetCookies(s.Domain, s.Cookies)
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
	headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	headers.Add("sentry-trace", "3e630a83fc154dcb8603694c5f14dabd-88912279c7b81e26-0")
}
