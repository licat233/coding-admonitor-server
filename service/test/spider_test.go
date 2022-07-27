package test

import (
	"admonitor-server/common/config"
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/utils"
	"admonitor-server/service/api_server/logic"
	"admonitor-server/service/monitor_server/core"
	"admonitor-server/service/monitor_server/spider_list/line_spider"
	"admonitor-server/service/monitor_server/spider_list/popin_spider"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gookit/goutil/dump"
	"github.com/spf13/viper"
)

func Test_LineSpider(t *testing.T) {
	config.Initialize()
	s := line_spider.NewLineSpider()
	s.InitAccounts()
}

func Test_PopinSpider(t *testing.T) {
	config.Initialize()
	s := popin_spider.NewPopinSpider()
	s.InitAccounts()
}

func Test_Popin(t *testing.T) {
	config.Initialize()
	s := popin_spider.NewPopinSpider()
	// s.GetCookies()
	timer := time.NewTicker(time.Second * 2)
	count := 0
	for count < 1 {
		res := s.Statistics(time.Now())
		if res != nil {
			fmt.Println("爬取成功")
		} else {
			fmt.Println("爬取失败")
		}
		count++
		<-timer.C
	}
}

func Test_PopinHeatbeat(t *testing.T) {
	config.Initialize()
	s := popin_spider.NewPopinSpider()
	s.Heartbeat()
}

func Test_PopinAccountInfo(t *testing.T) {
	config.Initialize()
	s := popin_spider.NewPopinSpider()
	s.Heartbeat()
}

func Test_nil(t *testing.T) {
	var list []string = nil
	for _, v := range list {
		fmt.Println(v)
	}
}

func Test_colly(t *testing.T) {
	myurl := "https://licat.coding.net/p/Admonitor"
	mycookies := []*http.Cookie{
		{
			Name:   "ses",
			Value:  "zbcdefg",
			Path:   "/p",
			Domain: ".coding.net",
		},
		{
			Name:   "currency",
			Value:  "TWD",
			Path:   "/p",
			Domain: ".coding.net",
		},
	}

	c := colly.NewCollector()
	c.SetCookies(myurl, mycookies)
	cCookies := c.Cookies(myurl)
	dump.Println(cCookies)
}

type MyCookieJar struct {
	Storage map[string][]*http.Cookie
}

func (m *MyCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if _, ok := m.Storage[u.Host]; !ok {
		m.Storage[u.Path] = cookies
	}
}

func (m *MyCookieJar) Cookies(u *url.URL) []*http.Cookie {
	if v, ok := m.Storage[u.Host]; ok {
		return v
	}
	return nil
}

func Test_colly2(t *testing.T) {
	requ := "http://www.baidu.com"
	u, _ := url.Parse(requ)
	cookies := []*http.Cookie{
		{
			Name:   "ses",
			Value:  "zbcdefg",
			Path:   "/platform",
			Domain: "ads.popin.cc",
		},
		{
			Name:   "currency",
			Value:  "TWD",
			Path:   "/platform",
			Domain: "ads.popin.cc",
		},
	}
	cj := &MyCookieJar{
		Storage: map[string][]*http.Cookie{},
	}
	cj.SetCookies(u, cookies)
	//
	c := colly.NewCollector()
	c.SetCookieJar(cj)
	cCookies := c.Cookies(requ)
	dump.Println(cookies)
	dump.Println(cCookies)
}

func Test_ReadLog(t *testing.T) {
	config.Initialize()
	lines, err := utils.ReadByLine(path.Join(utils.CurrentProjectRootPath(), viper.GetString("zap.log-dir"), viper.GetString("zap.log-file")))
	if err != nil {
		log.Fatalln(err)
		return
	}
	logData := []*logic.LogInfo{}
	for _, line := range lines {
		loginfo := &logic.LogInfo{}
		if err := json.Unmarshal(line, loginfo); err != nil {
			global.LOG.Error(errorx.New("反序列log信息失败", err))
			continue
		}
		logData = append(logData, loginfo)
	}
	dump.Println(logData)
}

func Test_core(t *testing.T) {
	config.Initialize()
	c := core.NewCoreServer()
	// c.CoreLogic.SpiderListRegister("line", line_spider.NewSpiderMgr)
	c.CoreLogic.SpiderListRegister("popin", popin_spider.NewSpiderMgr)
	//初始化spiders
	c.CoreLogic.InitSpiders()

	timer := time.NewTicker(time.Second * 2)
	count := 0
	for count < 1 {
		//测试
		c.CoreLogic.InsertNewestData()
		fmt.Println("爬取成功")
		count++
		<-timer.C
	}
}

func Test_Once(t *testing.T) {
	config.Initialize()
	c := core.NewCoreServer()
	c.CoreLogic.SpiderListRegister("line", line_spider.NewSpiderMgr)
	c.CoreLogic.SpiderListRegister("popin", popin_spider.NewSpiderMgr)
	c.CoreLogic.InitSpiders()
	c.Once()
}
