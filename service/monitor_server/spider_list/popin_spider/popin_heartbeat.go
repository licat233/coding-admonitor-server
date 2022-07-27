package popin_spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

/** Heartbeat 通过心跳检测cookie是否有效
 * 如果失效，则重新登录，更新cookie
 * 4 分钟内不能重复心跳
 */
func (s *PopinSpider) Heartbeat() {
	if s.Close {
		panic(fmt.Errorf("%s蜘蛛已经关闭", s.Channel.Name))
	}
	//不能心跳太快，4分内不能重复心跳
	if time.Since(s.LastHeartbeatTime) < 240*time.Second {
		return
	}
	reqUrl := "https://ads.popin.cc/platform/api/campaignManager/accountInfo"
	var e error
	c := colly.NewCollector()
	c.SetCookies(reqUrl, s.Cookies)
	c.OnRequest(func(req *colly.Request) {
		s.accountInfoHeader(req.Headers)
	})
	c.OnError(func(resp *colly.Response, err error) {
		e = s.OnErrorHandler(resp, err)
	})
	c.Visit(reqUrl) //这里的返回的error与onError里面的error一致
	c.Wait()
	if e != nil {
		s.SpiderErrorHandler(e, s.Heartbeat, s.GetCookies)
	}
	//访问成功，更新spider的cookies
	s.MergeCookies(c.Cookies(reqUrl))
	//更新到数据库
	if err := s.UpdateModelCookies(); err != nil {
		msg := errorx.New("mysql更新cookie数据失败", err)
		global.LOG.Error(msg)
	}
}
