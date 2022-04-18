package spider

import "net/http"

func userAgent() string {
	res := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
	return res
}

func headerHandler(headers *http.Header) {
	headers.Add("Accept", "application/json")
	headers.Add("User-Agent", userAgent())
	headers.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,ja;q=0.7")
	headers.Add("Connection", "keep-alive")
	// headers.Add("Cookie", cookieData())
	headers.Add("Referer", "https://ads.popin.cc/platform/new_dashboard")
	headers.Add("Referrer-Policy", "strict-origin-when-cross-origin")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "macOS")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Site", "same-origin")
}
