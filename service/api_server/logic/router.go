package logic

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func (a *ApiServer) Router() {
	api := a.Gin.Group("/api_server")
	{
		api.GET("/ping", Auth(), a.Pong)
		api.POST("/authorizetion", a.ReqAuthorizetion)
		api.POST("/statistics", Auth(), a.ReqStatistics)
		api.POST("/channelsAccounts", Auth(), a.ReqChannelsAndAccounts)
		api.POST("/channellist", Auth(), a.ReqChannelList)
		api.POST("/accountlist", Auth(), a.ReqAccountList)
		api.POST("/monitor_logs", Auth(), a.MonitorLogs)
		api.POST("/channel/status", Auth(), a.ModifyChannelStatus)
		api.POST("/account/status", Auth(), a.ModifyAccountStatus)
	}
	a.Gin.GET("/swagger/docs", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
