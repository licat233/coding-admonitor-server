package logic

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var jwtkey = []byte(viper.GetString("api-server.jwtkey"))

//SetToken 颁发token
func SetToken(openid, sessionKey string) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		OpenId:     openid,
		SessionKey: sessionKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "licat",      // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		msg := errorx.New("生成TokenString错误", err)
		global.LOG.Error(msg)
	}
	return tokenString, err
}

//getting 解析token
func Getting(ctx *gin.Context) *Claims {
	tokenString := ctx.GetHeader("Authorization")
	//vcalidate token formate
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		ctx.Abort()
		return nil
	}

	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		ctx.Abort()
		return nil
	}
	return claims
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
