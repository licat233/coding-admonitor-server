package popin_spider

import (
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

func statisticsURL(date time.Time) string {
	d := date.Format("20060102")
	return fmt.Sprint("https://ads.popin.cc/platform/api/campaignManager/dashboard/getAccountDashboardList?dateType=1&timezone=UTC%2B8&startDate=", d, "&endDate=", d, "&currency=USD&accountName=&accountIds=&currentPage=1&pageSize=50&sortField=&sortVal=")
}

//Statistics 获取统计数据
//如果 spider已经停止，则返回 nil，
//如果心跳失败，报401错误，则重新Login获取新的cookie
func (s *PopinSpider) Statistics(date time.Time) (resp *StatisticsResp) {
	if err := s.ValidateCookies(s.Cookies); err != nil {
		global.LOG.Error(err)
		panic(err)
	}

	if s.Close {
		panic(fmt.Errorf("%s蜘蛛已经关闭", s.Channel.Name))
	}
	reqUrl := statisticsURL(date)
	//短时间内请求，且url相同，不再请求，直接返回缓存数据
	if time.Since(s.LastStatisticsTime) < 30*time.Second && s.LastStatisticsUrl == reqUrl {
		// 判断是否为nil
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
		if resp.Code != 0 {
			e = errors.New(string(r.Body))
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
			//当前无权限的cookies
			// global.LOG.Info("spider中的cookies:", s.Cookies)
			// global.LOG.Info("colly中的cookies:", c.Cookies(reqUrl))
			//重新登录，再重新请求
			s.GetCookies()
			resp = s.Statistics(date)
		})
	}
	return resp
}
