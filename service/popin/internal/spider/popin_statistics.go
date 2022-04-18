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

/** StatisticsResp Statistics的响应 **/
type StatisticsResp struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

//Detail 广告数据详细
type Detail struct {
	AccountCurrency string `json:"accountCurrency"`
	AccountID       string `json:"accountId"`
	AccountName     string `json:"accountName"`
	Cpc             string `json:"avgCpc"`
	Click           string `json:"click"`
	Cv              string `json:"conversion"`
	Cpa             string `json:"conversionPrice"`
	Cvr             string `json:"conversionRate"`
	Ctr             string `json:"ctr"`
	Currency        string `json:"currency"`
	Imp             string `json:"impression"`
	Cost            string `json:"spend"`
}

//Summary 数据归总
type Summary struct {
	AvgCpc                 string `json:"avgCpc"`
	Click                  string `json:"click"`
	Conversion             string `json:"conversion"`
	ConversionCost         string `json:"conversionCost"`
	ConversionPrice        string `json:"conversionPrice"`
	ConversionRate         string `json:"conversionRate"`
	Ctr                    string `json:"ctr"`
	Cv                     string `json:"cv"`
	CvAddPaymentInfo       string `json:"cvAddPaymentInfo"`
	CvAddToCar             string `json:"cvAddToCar"`
	CvAddToWishlist        string `json:"cvAddToWishlist"`
	CvAppInstall           string `json:"cvAppInstall"`
	CvCompleteRegistration string `json:"cvCompleteRegistration"`
	CvLead                 string `json:"cvLead"`
	CvOther                string `json:"cvOther"`
	CvPurchase             string `json:"cvPurchase"`
	CvStartCheckout        string `json:"cvStartCheckout"`
	CvViewContent          string `json:"cvViewContent"`
	Impression             string `json:"impression"`
	Mcv                    string `json:"mcv"`
	Search                 string `json:"search"`
	Spend                  string `json:"spend"`
}

//Data 数据载体
type Data struct {
	CurrentPage int      `json:"currentPage"`
	Detail      []Detail `json:"detail"`
	Summary     Summary  `json:"summary"`
	Total       int      `json:"total"`
}

func satisticsURL(date time.Time) string {
	d := date.Format("20060102")
	return fmt.Sprint("https://ads.popin.cc/platform/api/campaignManager/dashboard/getAccountDashboardList?dateType=1&timezone=UTC%2B8&startDate=", d, "&endDate=", d, "&currency=USD&accountName=&accountIds=&currentPage=1&pageSize=50&sortField=&sortVal=")
}

//Statistics 获取统计数据
//如果 spider已经停止，则返回 nil，
//如果心跳失败，报401错误，则重新Login获取新的cookie
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
		headerHandler(req.Headers)
	})
	Colly.SetCookies(s.Domain, s.Cookies)
	var e error
	Colly.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = errors.New("响应状态码不是200")
			return
		}
		resp = new(StatisticsResp)
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
		global.LOG.Error("获取statistic数据无访问权限，准备重新登录")
		//调用heartbeat检验cookie的有效性
		s.Heartbeat()
		return
	}
	panic(fmt.Errorf(errorx.New("popin_statistics失败", err)))
}
