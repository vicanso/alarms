package main

import (
	"crypto/tls"
	"net/http"

	"github.com/vicanso/alarms/config"
	"github.com/vicanso/alarms/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gopkg.in/gomail.v2"
)

type AlarmParams struct {
	Service  string `json:"service,omitempty" valid:"runelength(1|30)"`
	Category string `json:"category,omitempty" valid:"runelength(1|30)"`
	Message  string `json:"message,omitempty" valid:"runelength(1|500)"`
	Token    string `json:"token,omitempty" valid:"runelength(1|30)"`
}

var (
	mailDialer *gomail.Dialer
	mailSender string
)

func init() {
	mailConfig := config.GetMailConfig()
	if mailConfig.Host != "" {
		mailSender = mailConfig.User
		mailDialer = gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.User, mailConfig.Password)
		mailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func main() {
	e := elton.New()

	// panic处理
	e.Use(middleware.NewRecover())

	// 出错处理
	e.Use(middleware.NewError(middleware.ErrorConfig{
		ResponseType: "json",
	}))

	// 默认的请求数据解析
	e.Use(middleware.NewDefaultBodyParser())

	// 响应数据转换为json
	e.Use(middleware.NewDefaultResponder())

	receivers := config.GetStringSlice("alarm.receiver")
	token := config.GetString("alarm.token")
	e.POST("/alarms", func(c *elton.Context) (err error) {
		params := AlarmParams{}
		err = validate.Do(&params, c.RequestBody)
		if err != nil {
			return
		}
		if params.Token != token {
			err = hes.New("token is invalid")
			return
		}
		m := gomail.NewMessage()
		m.SetHeader("From", mailSender)
		m.SetHeader("To", receivers...)
		m.SetHeader("Subject", params.Service+":"+params.Category)
		m.SetBody("text/plain", params.Message)
		err = mailDialer.DialAndSend(m)
		if err != nil {
			return
		}
		c.NoContent()
		return
	})

	// http1与http2均支持
	e.Server = &http.Server{
		Handler: h2c.NewHandler(e, &http2.Server{}),
	}

	err := e.ListenAndServe(config.GetListen())
	if err != nil {
		panic(err)
	}
}
