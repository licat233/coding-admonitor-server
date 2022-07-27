package popin_spider

import (
	"admonitor-server/service/monitor_server/core/spider"
	"net/http"
)

type LoginResp struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func (s *PopinSpider) addHeader(headers *http.Header) {
	headers.Add("authority", "ads.popin.cc")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("accept", "application/json, text/plain, */*")
	headers.Add("dnt", "1")
	headers.Add("content-type", "application/json; charset=UTF-8")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("user-agent", spider.UserAgent)
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("origin", "https://ads.popin.cc")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("referer", "https://ads.popin.cc/platform/user/login")
	headers.Add("accept-language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
	// headers.Add("cookie", s.HttpCookiesToString())
}

/** AccountInfoResp的响应 **/
type AccountInfoResp struct {
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

//UserTags
type UserTags struct {
	Label string `json:"label"`
	Uv    int    `json:"uv"`
	Value string `json:"value"`
}

func (s *PopinSpider) accountInfoHeader(headers *http.Header) {
	headers.Add("authority", "ads.popin.cc")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("accept", "application/json")
	headers.Add("dnt", "1")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("user-agent", spider.UserAgent)
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("referer", "https://ads.popin.cc/platform/new_dashboard")
	headers.Add("accept-language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
}

func statisticsHeader(headers *http.Header) {
	headers.Add("Accept", "application/json")
	headers.Add("User-Agent", spider.UserAgent)
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
	headers.Add("Connection", "keep-alive")
	headers.Add("Referer", "https://ads.popin.cc/platform/new_dashboard")
	headers.Add("Referrer-Policy", "strict-origin-when-cross-origin")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "macOS")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Site", "same-origin")
}
