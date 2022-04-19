package spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"errors"
	"fmt"
	"net/http"

	"github.com/gocolly/colly/v2"
)

/**
 * popin即使登录成功，需要验证过后，拿到最后一个关键cookie，才能正常使用
 */

/** authorize 验证身份，获取最终的cookie
 * 由于popin采用302链式跳转，所以这里采用回调方式
 * @param url type string 访问链接
 * @return statusCode type int 返回http状态码，当返回200时，表示验证成功
 * 注意：此函数只能由login()调用
 */
func (s *Spider) authorize() {
	url := "https://ads.popin.cc/platform/main/login?type=m&source=us"
	c := colly.NewCollector()
	c.OnRequest(func(req *colly.Request) {
		loginHeader(req.Headers)
	})
	c.SetCookies(url, s.Cookies)
	var e error
	var statusCode int
	c.OnResponse(func(r *colly.Response) {
		statusCode = r.StatusCode
	})
	c.OnError(func(resp *colly.Response, err error) {
		statusCode = resp.StatusCode
		e = err
	})
	if err := c.Visit(url); err != nil {
		s.authorizeError(err)
		return
	}
	c.Wait()
	if e != nil {
		s.authorizeError(e)
		return
	}
	if statusCode != 200 {
		s.authorizeError(fmt.Errorf("statusCode: %d", statusCode))
	}
	cookies := c.Cookies(url)
	if len(cookies) == 0 {
		s.authorizeError(errors.New("没有获取到cookie"))
	}
	//echo
	// for _, cookie := range cookies {
	// 	fmt.Printf("%s=%s\n", cookie.Name, cookie.Value)
	// }
	global.LOG.Info("popin鉴权成功")
	//将新的cookies与原cookies合并
	s.mergeCookies(cookies)
}

func (s *Spider) authorizeError(err error) {
	s.Close = true
	panic(fmt.Errorf(errorx.New("popin鉴权失败", err)))
}

func authorizeHeader(headers *http.Header) {
	headers.Add("authority", "ads.popin.cc")
	headers.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"macOS\"")
	headers.Add("upgrade-insecure-requests", "1")
	headers.Add("dnt", "1")
	headers.Add("user-agent", UserAgent)
	headers.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-user", "?1")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("referer", "https://ads.popin.cc/platform/user/login")
	headers.Add("accept-language", "zh-CN,zh;q=0.9")
}
