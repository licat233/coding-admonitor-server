package popin_spider

import (
	"admonitor-server/service/model"
	"admonitor-server/service/monitor_server/core/spider"
	"net/http"
	"time"
)

// PopinSpider 注意：该结构体务必实现 spider.SpiderMgr 接口
type PopinSpider struct {
	LastHeartbeatTime time.Time
	*spider.Spider
}

// 检测PopinSpider是否实现了spider.SpiderMgr接口
// 1） _为了避免变量未使用编译的时候报错
// 2）_的类型为spider.SpiderMgr，接口的值为PopinSpider的地址，(nil)表示该地址为nil。
var _ spider.SpiderMgr = (*PopinSpider)(nil)

func NewPopinSpider() *PopinSpider {
	domain := "ads.popin.cc"
	baseCookies := []*http.Cookie{
		spider.GenHttpCookie("currency", "TWD", domain),
		spider.GenHttpCookie("new_db_trend", "spend", domain),
		spider.GenHttpCookie("mgruli", "n%3D%26t%3Dmsnpopin", domain),
		spider.GenHttpCookie("mgvdata", "uYK9xpYQDCnhzRv9", domain),
		spider.GenHttpCookie("db_campaign", "Management#campaignName%2CCost#spend%2CCampaigns#campaignBalance%2CTrafficAcquisition#cpc%2CTrafficAcquisition#impression%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr%2Cconversion#conversion%2Cconversion#conversionRate", domain),
		spider.GenHttpCookie("db_account", "Management#accountName%2CCost#spend%2CTrafficAcquisition#avgCpc%2CTrafficAcquisition#impression%2CTrafficAcquisition#click%2CTrafficAcquisition#ctr%2Cconversion#conversion%2Cconversion#conversionRate%2Cconversion#conversionPrice", domain),
	}
	s := spider.NewSpider(&model.Channel{
		Name:   "popin",
		Domain: domain,
	}, baseCookies...)
	return &PopinSpider{
		LastHeartbeatTime: time.Time{},
		Spider:            s,
	}
}

//新建popin实现的spider.SpiderMgr接口，用于注册
func NewSpiderMgr() spider.SpiderMgr {
	return NewPopinSpider()
}
