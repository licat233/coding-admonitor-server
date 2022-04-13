package errorx

import "fmt"

type Errorx struct {
	Message string //错误消息
	Addr    string //错误地址
	Error   error  //具体错误
}

func New(msg string, err error) string {
	return fmt.Sprintf("【错误信息:%s】【错误原因:%s】", msg, err)
}
