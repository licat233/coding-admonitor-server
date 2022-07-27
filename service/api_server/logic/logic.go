package logic

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"admonitor-server/common/respx"
	"admonitor-server/common/utils"
	"admonitor-server/service/model"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApiServer struct {
	ChannelModel    *model.ChannelModel
	AccountModel    *model.AccountModel
	StatisticsModel *model.StatisticsModel
	Gin             *gin.Engine
}

func NewApiServer() *ApiServer {
	return &ApiServer{
		ChannelModel:    model.NewChannelModel(),
		AccountModel:    model.NewAccountModel(),
		StatisticsModel: model.NewStatisticsModel(),
		Gin:             gin.Default(),
	}
}

func (a *ApiServer) Start() {
	a.Router()
	port := viper.GetString("api-server.port")
	addr := fmt.Sprintf(":%s", port)
	if err := a.Gin.Run(addr); err != nil {
		global.LOG.Error("gin启动失败:", zap.Error(err))
	}
}

func (a *ApiServer) ReqStatistics(c *gin.Context) {
	req := new(StatisticsRequest)
	err := c.BindJSON(req)
	if err != nil {
		errorx.ExternalError(c, err)
		return
	}
	startDate, err := utils.ParseDate(req.StartDate)
	if err != nil {
		errorx.ExternalError(c, errors.New("'startDate' field value error"))
		return
	}
	endDate, err := utils.ParseDate(req.EndDate)
	if err != nil {
		errorx.ExternalError(c, errors.New("'endDate' field value error"))
		return
	}
	if startDate.After(endDate) {
		errorx.ExternalError(c, errors.New("开始时间大于结束时间"))
		return
	}
	now := time.Now().Local()
	if startDate.After(now) {
		errorx.ExternalError(c, errors.New("开始时间大于当前时间"))
		return
	}
	statisticsType := req.Type
	if statisticsType != model.LOGIC && statisticsType != model.NEWEST && statisticsType != model.DAY {
		errorx.ExternalError(c, errors.New("'type' field value error"))
	}
	lists, err := a.StatisticsModel.FindsByAccountIds(req.ChannelId, req.AccountIds, startDate, endDate, statisticsType)
	if len(lists) == 0 || errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, respx.SingleResp("查询结果为空", utils.StatisticConvert(utils.CalculateTotalStatistics(lists))))
		return
	} else if err != nil {
		errorx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, respx.SingleResp("查询成功", utils.StatisticConvert(utils.CalculateTotalStatistics(lists))))
}

func (a *ApiServer) ReqChannelsAndAccounts(c *gin.Context) {
	req := new(ChannelsAndAccountsRequest)
	err := c.BindJSON(req)
	if err != nil {
		errorx.ExternalError(c, err)
		return
	}
	if req.ChannelStatus != -1 && req.ChannelStatus != 0 && req.ChannelStatus != 1 {
		errorx.ExternalError(c, errors.New("'channelStatus' param error"))
		return
	}
	if req.AccountStatus != -1 && req.AccountStatus != 0 && req.AccountStatus != 1 {
		errorx.ExternalError(c, errors.New("'accountStatus' param error"))
		return
	}
	var channels []*model.Channel
	switch req.ChannelStatus {
	case 0:
		channels, err = a.ChannelModel.FindAllByStatus(false)
	case 1:
		channels, err = a.ChannelModel.FindAllByStatus(true)
	case -1:
		channels, err = a.ChannelModel.FindAll()
	default:
		errorx.ExternalError(c, errors.New("'channelStatus' param error"))
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, respx.SingleResp("查询结果为空", nil))
		return
	} else if err != nil {
		errorx.InternalError(c, err)
		return
	}
	allAccount, err := a.AccountModel.FindAll()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errorx.InternalError(c, errors.New("数据库无任何账号信息"))
		return
	} else if err != nil {
		errorx.InternalError(c, err)
		return
	}
	resp := []*ChannelsResponse{}
	for _, Channel := range channels {
		Channel.Password = "" //去除密码
		accountList := AccountsFilter(allAccount, Channel.ID, req.AccountStatus)
		Accounts := []*AccountInfo{}
		hasCheck := false
		for _, Account := range accountList {
			check := false
			//默认选择第一个status为true的账号
			if !hasCheck && Account.Status {
				hasCheck = true
				check = true
			}
			// if index == 0 { //默认选中第一个
			// 	check = true
			// }
			accountInfo := &AccountInfo{
				Id:        Account.ID,
				Account:   Account,
				IsChecked: check,
				Disabled:  false,
			}
			Accounts = append(Accounts, accountInfo)
		}
		if !hasCheck && len(Accounts) > 0 {
			Accounts[0].IsChecked = true
		}
		resp = append(resp, &ChannelsResponse{
			Id:       Channel.ID,
			Channel:  Channel,
			Accounts: Accounts,
		})
	}
	c.JSON(http.StatusOK, respx.SingleResp("查询成功", resp))
}

func (a *ApiServer) GetChannelList(c *gin.Context, statusCode int) (resp []*ChannelListResponse, err error) {
	var channels []*model.Channel
	switch statusCode {
	case 0:
		channels, err = a.ChannelModel.FindAllByStatus(false)
	case 1:
		channels, err = a.ChannelModel.FindAllByStatus(true)
	case -1:
		channels, err = a.ChannelModel.FindAll()
	default:
		errorx.ExternalError(c, errors.New("channel 'status' param error"))
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, respx.SingleResp("查询结果为空", nil))
		err = nil
		return
	} else if err != nil {
		errorx.InternalError(c, err)
		return
	}
	for _, channel := range channels {
		channel.Password = "" //去除密码
		resp = append(resp, &ChannelListResponse{
			Id:     channel.ID,
			Detail: channel,
		})
	}
	return
}

func (a *ApiServer) ReqChannelList(c *gin.Context) {
	p := new(StatusRequest)
	err := c.BindJSON(p)
	if err != nil {
		errorx.ExternalError(c, err)
		return
	}
	resp, err := a.GetChannelList(c, p.Status)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, respx.SingleResp("查询成功", resp))
}

func (a *ApiServer) GetAccountList(c *gin.Context, statusCode int) (resp []*AccountListResponse, err error) {
	var accounts []*model.Account
	switch statusCode {
	case 0:
		accounts, err = a.AccountModel.FindAllByStatus(false)
	case 1:
		accounts, err = a.AccountModel.FindAllByStatus(true)
	case -1:
		accounts, err = a.AccountModel.FindAll()
	default:
		errorx.ExternalError(c, errors.New("account 'status' param error"))
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, respx.SingleResp("查询结果为空", nil))
		err = nil
		return
	} else if err != nil {
		errorx.InternalError(c, err)
		return
	}
	for _, account := range accounts {
		resp = append(resp, &AccountListResponse{
			Id:     account.ID,
			Detail: account,
		})
	}
	return
}

func (a *ApiServer) ReqAccountList(c *gin.Context) {
	p := new(StatusRequest)
	err := c.BindJSON(p)
	if err != nil {
		errorx.ExternalError(c, err)
		return
	}
	resp, err := a.GetAccountList(c, p.Status)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, respx.SingleResp("查询成功", resp))
}

func (a *ApiServer) Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "pong")
}

func (a *ApiServer) MonitorLogs(ctx *gin.Context) {
	lines, err := utils.ReadByLine(path.Join(utils.CurrentProjectRootPath(), viper.GetString("zap.log-dir"), viper.GetString("zap.log-file")))
	if err != nil {
		errorx.InternalError(ctx, err)
		return
	}
	logs := []string{}
	for _, line := range lines {
		logs = append(logs, string(line))
	}
	ctx.JSON(http.StatusOK, respx.SingleResp("请求成功", logs))
}

//ModifyChannelStatus 修改渠道状态
func (a *ApiServer) ModifyChannelStatus(ctx *gin.Context) {
	p := new(ModifyStatusRequest)
	err := ctx.BindJSON(p)
	if err != nil {
		errorx.ExternalError(ctx, err)
		return
	}
	err = a.ChannelModel.UpdateStatus(p.Ids, p.Status)
	if err != nil {
		errorx.InternalError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, respx.StateResp("修改成功"))
}

//ModifyAccountStatus 修改账号状态
func (a *ApiServer) ModifyAccountStatus(ctx *gin.Context) {
	p := new(ModifyStatusRequest)
	err := ctx.BindJSON(p)
	if err != nil {
		errorx.ExternalError(ctx, err)
		return
	}
	err = a.AccountModel.UpdateStatus(p.Ids, p.Status)
	if err != nil {
		errorx.InternalError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, respx.StateResp("修改成功"))
}
