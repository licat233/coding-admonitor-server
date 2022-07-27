package line_spider

import (
	"admonitor-server/service/monitor_server/core/spider"
	"fmt"
	"net/http"
	"time"
)

type LoginResp struct {
	Status string `json:"status"`
}

//loginFaild statusCode = 400
// type loginFaild struct {
// 	Status       string         `json:"status"` //failed
// 	Notices      []interface{}  `json:"notices"`
// 	GlobalErrors []GlobalErrors `json:"globalErrors"`
// 	FieldErrors  []interface{}  `json:"fieldErrors"`
// }

// type GlobalErrors struct {
// 	Message string `json:"message"`
// }

//stutusCode = 200，需要验证码，暂时无法破解验证码
// type captachResp struct {
// 	Status          string         `json:"status"` //needCaptchaVerification
// 	CaptchaKey      string         `json:"captchaKey"`
// 	CaptchaImageURI string         `json:"captchaImageUri"`
// 	GlobalErrors    []GlobalErrors `json:"globalErrors"`
// }

func loginPageHeader(headers *http.Header) {
	headers.Add("Connection", "keep-alive")
	headers.Add("Cache-Control", "max-age=0")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("DNT", "1")
	headers.Add("Upgrade-Insecure-Requests", "1")
	headers.Add("User-Agent", spider.UserAgent)
	headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	headers.Add("Sec-Fetch-Site", "none")
	headers.Add("Sec-Fetch-Mode", "navigate")
	headers.Add("Sec-Fetch-User", "?1")
	headers.Add("Sec-Fetch-Dest", "document")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
}

func loginHeader(headers *http.Header) {
	headers.Add("Connection", "keep-alive")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("DNT", "1")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("User-Agent", spider.UserAgent)
	headers.Add("Content-Type", "application/json;charset=UTF-8")
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("Origin", "https://account.line.biz")
	headers.Add("Sec-Fetch-Site", "same-origin")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://account.line.biz/login?redirectUri=https%3A%2F%2Fadmanager.line.biz%2Flogin%2F&scope=email")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
}

func accountInfoHeader(headers *http.Header) {
	headers.Add("Connection", "keep-alive")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("DNT", "1")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("User-Agent", spider.UserAgent)
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("Sec-Fetch-Site", "same-origin")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://admanager.line.biz/adaccount/")
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
}

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
	headers.Add("User-Agent", spider.UserAgent)
	headers.Add("sentry-trace", "3e630a83fc154dcb8603694c5f14dabd-88912279c7b81e26-0")
}
