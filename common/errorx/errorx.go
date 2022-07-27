package errorx

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Errorx struct {
	Message string //错误消息
	Addr    string //错误地址
	Error   error  //具体错误
}

func New(msg string, err error) string {
	return fmt.Sprintf("【错误信息:%s】【错误原因:%s】", msg, err)
}

func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  false,
		"code":    http.StatusInternalServerError,
		"message": New("服务器内部错误", err),
	})
}

func ExternalError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  false,
		"code":    http.StatusBadRequest,
		"message": New("请求参数错误", err),
	})
}
