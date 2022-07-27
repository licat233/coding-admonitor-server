package line_spider

import (
	"encoding/json"
	"errors"
	"fmt"
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

func statisticsURL(date time.Time) string {
	todayDate := date.Format("2006-01-02")
	return fmt.Sprintf("https://admanager.line.biz/api/adservice/v1/my/adaccounts/statistics?groupId=&page=1&size=100&searchKey=&sort=&reportSince=%s&reportUntil=%s", todayDate, todayDate)
}

//Statistics 获取统计数据，爬虫的核心方法
//如果 spider已经停止，则返回 nil，
//如果报401错误，则重新Login获取新的cookie
func (s *LineSpider) Statistics(date time.Time) (resp *StatisticsResp) {
	//line渠道不需要验证cookies，默认cookies为nil
	// if err := s.ValidateCookies(s.Cookies); err != nil {
	// 	global.LOG.Error(err)
	// 	panic(err)
	// }

	if s.Close {
		panic(fmt.Errorf("%s蜘蛛已经关闭", s.Channel.Name))
	}

	reqUrl := statisticsURL(date)
	//短时间内请求，且url相同，不再请求，直接返回缓存数据
	if time.Since(s.LastStatisticsTime) < 30*time.Second && s.LastStatisticsUrl == reqUrl {
		//判断是否为nil
		if s.CacheStatisticsData != nil {
			//类型断言
			if v, ok := s.CacheStatisticsData.(*StatisticsResp); ok {
				return v
			}
		}
	}

	c := colly.NewCollector()
	c.SetCookies(reqUrl, s.Cookies)
	c.OnRequest(func(req *colly.Request) {
		statisticsHeader(req.Headers)
	})
	var e error
	c.OnResponse(func(r *colly.Response) {
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
		s.CacheStatisticsData = resp
		s.LastStatisticsTime = time.Now().Local()
		s.LastStatisticsUrl = reqUrl
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = s.OnErrorHandler(resp, err)
	})
	c.Visit(reqUrl) //这里的返回的error与onError里面的error一致
	c.Wait()
	if e != nil {
		s.SpiderErrorHandler(e, func() {
			resp = s.Statistics(date)
		}, func() {
			//重新登录，再重新请求
			s.GetCookies()
			resp = s.Statistics(date)
		})
	}
	return resp
}
