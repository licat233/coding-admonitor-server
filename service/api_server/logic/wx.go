package logic

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/respx"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
)

func (a *ApiServer) ReqAuthorizetion(c *gin.Context) {
	p := new(AuthorizetionRequest)
	err := c.BindJSON(p)
	if err != nil {
		errorx.ExternalError(c, err)
		return
	}
	appid := viper.GetString("wxapp.appid")
	secret := viper.GetString("wxapp.secret")
	js_code := p.Code
	grant_type := viper.GetString("wxapp.grant_type")

	reqUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s", appid, secret, js_code, grant_type)
	s := colly.NewCollector()
	data := new(WeixinAuthResponese)
	s.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, data)
	})
	if e := s.Visit(reqUrl); e != nil || err != nil || len(data.OpenId) == 0 || len(data.SessionKey) == 0 {
		errorx.ExternalError(c, errors.New("验证未通过"))
		return
	}
	//获取到 openid 与 session_key以后，保存起来，生成 token
	tokenString, err := SetToken(data.OpenId, data.SessionKey)
	if err != nil {
		errorx.InternalError(c, errors.New("服务器错误"))
		return
	}
	c.JSON(http.StatusOK, respx.SingleResp("登录成功", &AuthorizetionResponse{Token: tokenString}))
}
