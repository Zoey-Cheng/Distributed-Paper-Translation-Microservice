package email

import (
	"net/smtp"

	"github.com/jordan-wright/email" // 邮件发送库
	"go-micro.dev/v4/config"         // 配置管理器
)

/**
* 根据配置创建邮件发送池
* @param config - 配置管理器
* @return email.Pool - 邮件发送池
 */
func NewEmailPool(config config.Config) *email.Pool {

	address := config.Get("email", "address").String("smtp.163.com:25")         // 获取地址
	username := config.Get("email", "username").String("golangproject@163.com") // 获取用户名
	password := config.Get("email", "password").String("GGFWFZWFQRDAKVYV")      // 获取密码
	host := config.Get("email", "host").String("smtp.163.com")                  // 获取SMTP主机

	// 创建邮件池
	pool, err := email.NewPool(
		address,
		4,
		smtp.PlainAuth(
			"",
			username,
			password,
			host,
		),
	)

	if err != nil { // 如果创建失败,panic
		panic(err)
	}

	return pool
}
