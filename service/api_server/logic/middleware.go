package logic

import "github.com/gin-gonic/gin"

//Auth 验证token
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if Getting(ctx) == nil {
			return
		}
		ctx.Next()
	}
}
