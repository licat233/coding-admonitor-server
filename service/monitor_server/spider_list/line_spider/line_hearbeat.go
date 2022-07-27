package line_spider

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gocolly/colly/v2"
)

func (s *LineSpider) Heartbeat() {}

//AccountInfo 获取账号信息
func (s *LineSpider) AccountInfo() {
	if s.Close {
		panic(fmt.Errorf("%s蜘蛛已经关闭", s.Channel.Name))
	}
	reqUrl := "https://ads.popin.cc/platform/api/campaignManager/accountInfo"
	var e error
	c := colly.NewCollector()
	c.SetCookies(reqUrl, s.Cookies)
	c.OnRequest(func(req *colly.Request) {
		accountInfoHeader(req.Headers)
	})
	resp := new(accountInfoResp)
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e = errors.New("响应状态码不是200")
			return
		}
		//反序列化
		err := json.Unmarshal(r.Body, resp)
		if err != nil {
			e = errors.New("json反序列化失败")
			return
		}
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = s.OnErrorHandler(resp, err)
	})
	c.Visit(reqUrl) //这里的返回的error与onError里面的error一致
	if e != nil {
		s.SpiderErrorHandler(e, s.AccountInfo, s.GetCookies)
	}
}
