package emailx

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/logx"
	"crypto/tls"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type Email struct {
	Username string
	Password string
	Host     string
	Port     int
	Message  *gomail.Message
	Dialer   *gomail.Dialer
	Receiver []string
}

var lastSendTime time.Time
var lastMsg string

func Send(msg string, to ...string) error {
	//【1分钟内】不允许重复发送【相同的信息】
	if lastMsg == msg && time.Since(lastSendTime) < time.Minute {
		return nil
	}
	lastMsg = msg
	lastSendTime = time.Now()
	email := &Email{
		Username: viper.GetString("email.username"),
		Password: viper.GetString("email.password"),
		Host:     viper.GetString("email.host"),
		Port:     viper.GetInt("email.port"),
		Receiver: viper.GetStringSlice("email.receiver"),
	}
	m := gomail.NewMessage()
	m.SetHeader("From", email.Username)
	m.SetHeader("To", email.Receiver...)
	m.SetHeader("Subject", "AdMonito通知")
	email.Message = m
	// 新建拨号器
	email.Dialer = gomail.NewDialer(email.Host, email.Port, email.Username, email.Password)
	// 关闭SSL协议认证
	email.Dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if len(to) != 0 {
		email.Receiver = append(email.Receiver, to...)
		email.Message.SetHeader("To", email.Receiver...)
	}
	email.Message.SetBody("text/html", emailBody(msg))
	if err := email.Dialer.DialAndSend(email.Message); err != nil {
		logx.NewLogger().Error(errorx.New("邮件发送失败", err))
		return err
	}
	return nil
}

func emailBody(msg string) string {
	index := strings.Index(msg, ":")
	title := msg
	content := ""
	if index != -1 {
		title = msg[:index]
		content = msg[index+1:]
	}
	return `<div style="max-width:600px;min-width:320px;margin:0 auto;background:#fff;border-collapse:collapse;font-size:18px">
	<div style="width:100%">
		<div style=background:#000;padding:30px>
			<img src=https://img.alicdn.com/imgextra/i2/917298378/O1CN01qBKkS22BlAzXe0eTI_!!917298378.png
				style=user-select:none;touch-action:none;pointer-events:none;max-width:100%>
		</div>
	</div>
	<div style="width:100%">
		<div style="padding:20px 15px">
			<h1 style="color:#000;border-bottom:1px solid;text-align:center;margin-bottom:30px">通&emsp;知</h1>
			<h2 style="color:red;font-size: 20px;">事件： ` + title + `</h2>
			<h3 style="color:#000;font-size: 18px;">原因：` + content + `</h3>
		</div>
	</div>
	<div style="position: fixed;left:0;bottom:0;width:100%">
		<div style=background:#eee;padding:10px;font-size:10px;text-align:center>AdMonitor系统</div>
	</div>
</div>`
}
