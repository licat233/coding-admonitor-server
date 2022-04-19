package emailx

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/logx"
	"crypto/tls"

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

func NewEmail() *Email {
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
	return email
}

func (e *Email) Send(msg string, to ...string) error {
	if len(to) != 0 {
		e.Receiver = append(e.Receiver, to...)
		e.Message.SetHeader("To", e.Receiver...)
	}
	e.Message.SetBody("text/html", emailBody(msg))
	if err := e.Dialer.DialAndSend(e.Message); err != nil {
		logx.NewLogger().Error(errorx.New("邮件发送失败", err))
		return err
	}
	return nil
}

func emailBody(msg string) string {
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
			<p style="color:red;font-size: 20px;">` + msg + `</p>
		</div>
	</div>
	<div style="position: fixed;left:0;bottom:0;width:100%">
		<div style=background:#eee;padding:10px;font-size:10px;text-align:center>AdMonitor系统</div>
	</div>
</div>`
}
