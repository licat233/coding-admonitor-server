package spider

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/service/model"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
)

const MaxLoginAcount int = 3 //最大登录次数，默认为 3 次，如果重复登录达到3次，就抛出错误，终止爬虫

//Spider 核心蜘蛛结构体，后期扩展的channelSpider可与其组合在一起
type Spider struct {
	Close               bool               //状态，注意：只能由login函数进行关闭
	Channel             *model.Channel     //当前spider的渠道record，每个渠道可以有多个账号,这里的渠道也代表登录账号，name并不是唯一的
	Cookies             []*http.Cookie     //存储cookie
	CacheStatisticsData any                //上次相同url请求结果缓存
	LastStatisticsTime  time.Time          //上次爬取statistics的时间
	LastStatisticsUrl   string             //上次爬取statistics的url
	LastLoginTime       time.Time          //上次登录时间，防止短时间重复登录
	CookieModel         *model.CookieModel //cookieModel
	LoginAcount         int                //登录次数
}

func NewSpider(setChannel *model.Channel, cookies ...*http.Cookie) *Spider {
	channel, err := setChannelModelRecord(setChannel)
	if err != nil {
		panic(errorx.New("创建Spider失败", err))
	}

	s := &Spider{
		Close:               false,
		Channel:             channel,
		Cookies:             cookies,
		CacheStatisticsData: nil,
		LastStatisticsTime:  time.Time{},
		LastStatisticsUrl:   "",
		LastLoginTime:       time.Time{},
		CookieModel:         model.NewCookieModel(),
		LoginAcount:         0,
	}
	s.InitSpiderCookies()
	return s
}

//setChannelModelRecord 设置channel，如果不存在，则创建新的，会panic
func setChannelModelRecord(setChannel *model.Channel) (*model.Channel, error) {
	channelName := setChannel.Name
	if len(channelName) == 0 {
		panic("请务必设置好渠道名")
	}
	channeModel := model.NewChannelModel()
	//查找是否存在指定名称的渠道
	channel, err := channeModel.FindByName(channelName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//不存在，那就创建一个新的，并且提示手动去数据库完善数据
		err := channeModel.InsertOne(setChannel)
		if err != nil {
			return nil, err
		}
		msg := fmt.Sprintf("【%s】渠道已经创建成功，请手动完善渠道信息", channelName)
		global.LOG.Info(msg)
		panic(msg)
	} else if err != nil {
		global.LOG.Error(errorx.New("mysql查询channel数据失败", err))
	}
	//如果存在了，就返回channel数据
	return channel, nil
}

//AllowLogin 在调用爬虫登录获取cookie前，判断是否允许
func (s *Spider) AllowLogin() bool {
	if s.Close {
		return false
	}
	//登录次数超过最大次数
	// if s.LoginAcount >= MaxLoginAcount {
	// 	msg := fmt.Sprintf("[%s]登录次数超过最大次数[%d]，已将爬虫关闭，具体原因请查看日志文件", s.Channel.Name, MaxLoginAcount)
	// 	s.Close = true
	// 	global.LOG.Error(msg)
	// 	panic(errors.New(msg))
	// }
	//30秒内不能重复登录
	if time.Since(s.LastLoginTime) < 30*time.Second {
		global.LOG.Error("请勿在30秒内重复登录")
		return false
	}
	//如果离上一次登录已经过去1个小时了，重置登录次数
	if time.Since(s.LastLoginTime) > 60*time.Minute {
		s.LoginAcount = 0
	}
	//允许登录
	//登录次数+1
	s.LoginAcount++
	//重置登录时间
	s.LastLoginTime = time.Now().Local()
	return true
}

//OnErrorHandler 将spider爬虫错误转化处理，转化为自定义的http错误，便于后期SpiderErrorHandler判断处理
func (s *Spider) OnErrorHandler(resp *colly.Response, err error) error {
	code := resp.StatusCode
	if httpError, ok := HttpErrors[code]; ok {
		return httpError
	}
	return err
}

//SpiderErrorHandler 爬虫错误处理
func (s *Spider) SpiderErrorHandler(err error, timeoutHandler func(), unauthHandler func()) {
	errmsg := err.Error()
	if strings.Contains(errmsg, "Client.Timeout") {
		//客户端超时，网络问题，再请求一次
		<-time.NewTimer(30 * time.Second).C //防止短时间内发生大量请求，我们30秒以后再请求一次
		timeoutHandler()
		return
	}

	switch err {
	case nil:
		return
	case HttpErrors[401]:
		global.LOG.Errorf("%s爬取无权限，将重新登录", s.Channel.Name)
		unauthHandler()
	case HttpErrors[408]:
		global.LOG.Errorf("%s爬取超时，将在30s后重新发起请求", s.Channel.Name)
		<-time.NewTimer(30 * time.Second).C //防止短时间内发生大量请求，我们30秒以后再请求一次
		timeoutHandler()
	default:
		//其它错误，关闭蜘蛛，抛出异常
		s.Close = true
		msg := fmt.Sprintf("[%s]爬取遇到未知错误，将关闭爬虫", s.Channel.Name)
		alert := errorx.New(msg, err)
		global.LOG.Error(alert)
		panic(errors.New(alert))
	}
}
