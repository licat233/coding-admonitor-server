package logic

import (
	"admonitor-server/service/model"

	"github.com/dgrijalva/jwt-go"
)

//Response 返回规范
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type StatusRequest struct {
	Status int `json:"status"` //状态，0关闭，1开启，-1表示无限制
}

//修改状态的请求
type ModifyStatusRequest struct {
	Ids    []int64 `json:"ids"`    //ids
	Status bool    `json:"status"` //状态
}

//获取日期范围内的广告监测数据
type StatisticsRequest struct {
	ChannelId  int64   `json:"channelId"`  //渠道id
	AccountIds []int64 `json:"accountIds"` //账号id
	StartDate  string  `json:"startDate"`  //开始日期
	EndDate    string  `json:"endDate"`    //结束日期
	Type       string  `json:"type"`       //数据类型
}

//GetChannels
type ChannelsAndAccountsRequest struct {
	ChannelStatus int `json:"channelStatus"` //状态，0关闭，1开启，-1表示无限制
	AccountStatus int `json:"accountStatus"` //状态，0关闭，1开启，-1表示无限制
}
type ChannelsResponse struct {
	Id       int64          `json:"id"`
	Channel  *model.Channel `json:"channelDetail"`
	Accounts []*AccountInfo `json:"accounts"`
}

type AccountInfo struct {
	Id        int64          `json:"id"`
	Account   *model.Account `json:"accountDetail"`
	IsChecked bool           `json:"isChecked"`
	Disabled  bool           `json:"disabled"`
}

//GetChannelList
type ChannelListResponse struct {
	Id     int64          `json:"id"`
	Detail *model.Channel `json:"channelDetail"`
}

//GetAccountList
type AccountListResponse struct {
	Id     int64          `json:"id"`
	Detail *model.Account `json:"accountDetail"`
}

//authorizetion
type AuthorizetionRequest struct {
	Code string `json:"code"` //小程序传过来的code
}

type AuthorizetionResponse struct {
	Token string `json:"token"`
}

type WeixinAuthResponese struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	// Errcode    int    `json:"errcode"`
	// Errmsg     string `json:"errmsg"`
}

//jwt-token
type Claims struct {
	OpenId     string
	SessionKey string
	jwt.StandardClaims
}

//logs
type LogInfo struct {
	Level    string `json:"level"`
	Datetime string `json:"ts"`
	Caller   string `json:"caller"`
	Msg      string `json:"msg"`
}
