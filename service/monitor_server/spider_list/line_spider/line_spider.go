package line_spider

import (
	"admonitor-server/service/model"
	"admonitor-server/service/monitor_server/core/spider"
)

//LineSpider 注意：该结构体务必实现 spider.SpiderMgr 接口
type LineSpider struct {
	*spider.Spider
}

// 检测LineSpider是否实现了spider.SpiderMgr接口
// 1） _为了避免变量未使用编译的时候报错
// 2）_的类型为spider.SpiderMgr，接口的值为LineSpider的地址，(nil)表示该地址为nil。
var _ spider.SpiderMgr = (*LineSpider)(nil)

func NewLineSpider() *LineSpider {
	s := spider.NewSpider(&model.Channel{
		Name:   "line",
		Domain: ".line.biz",
	})
	return &LineSpider{
		Spider: s,
	}
}

//新建lien实现的spider.SpiderMgr接口，用于注册
func NewSpiderMgr() spider.SpiderMgr {
	return NewLineSpider()
}
