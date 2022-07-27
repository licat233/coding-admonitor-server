package spider

import (
	"admonitor-server/service/model"
	"errors"
	"time"
)

var (
	ErrAuth    = errors.New("请求访问无权限")
	ErrTimeout = errors.New("请求超时")
	UserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
)

/**
 * 后期扩展的channelSpider结构体这样写
 * type ChannelSpider struct {
 *	*Spider //相当于继承Spider
 * }
 * ChannelSpider 需要实现 SpiderMgr 接口
 */

//SpiderMgr 蜘蛛管理员
//方便拓展，各个渠道只需实现该接口即可，务必实现该接口
type SpiderMgr interface {
	Info() *Spider //蜘蛛信息
	Heartbeat()    //心跳处理，内部可能会引发panic，请做好recover处理
	//爬取数据并且返回statistics数据，内部可能会引发panic，请做好recover处理
	GetStatistics(date time.Time, Type string, accounts []*model.Account) []*model.Statistics
	//获取网络cookies
	GetCookies()
}

var HttpErrors = map[int]error{
	100: errors.New("continue"),
	101: errors.New("witching protocols"),
	200: errors.New("ok"),
	201: errors.New("created"),
	202: errors.New("accepted"),
	203: errors.New("non-authoritative information"),
	204: errors.New("no content"),
	205: errors.New("reset content"),
	206: errors.New("partial content"),
	300: errors.New("multiple choices"),
	301: errors.New("moved permanently"),
	302: errors.New("found"),
	303: errors.New("see other"),
	304: errors.New("not modified"),
	305: errors.New("use proxy"),
	307: errors.New("temporary redirect"),
	400: errors.New("bad request"),
	401: errors.New("unauthorized"),
	402: errors.New("payment required"),
	403: errors.New("forbidden"),
	404: errors.New("not found"),
	405: errors.New("method not allowed"),
	406: errors.New("not acceptable"),
	407: errors.New("proxy authentication required"),
	408: errors.New("request time-out"),
	409: errors.New("conflict"),
	410: errors.New("gone"),
	411: errors.New("length required"),
	412: errors.New("precondition failed"),
	413: errors.New("request entity too large"),
	414: errors.New("request-uri too large"),
	415: errors.New("unsupported media type"),
	416: errors.New("requested range not satisfiable"),
	417: errors.New("expectation failed"),
	500: errors.New("internal server error"),
	501: errors.New("not implemented"),
	502: errors.New("bad gateway"),
	503: errors.New("service unavailable"),
	504: errors.New("gateway time-out"),
	505: errors.New("http version not supported"),
}
