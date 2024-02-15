package email

import (
	"net/smtp" // SMTP发送协议
	"testing"
	"time"

	"github.com/jordan-wright/email"     // 邮件发送库
	"github.com/stretchr/testify/assert" // 单元测试assert库
)

// 测试发送邮件
func TestMail(t *testing.T) {

	// 创建邮件池
	pool, err := email.NewPool(
		"smtp.163.com:25", // SMTP服务器地址
		4,                 // 最大连接数
		smtp.PlainAuth(
			"",                      // 用户名
			"golangproject@163.com", // 邮箱账号
			"GGFWFZWFQRDAKVYV",      // 授权码
			"smtp.163.com",          // SMTP服务器地址
		),
	)

	// 断言池创建是否错误
	assert.NoError(t, err)

	// 构建邮件内容
	em := email.NewEmail()
	em.From = ""
	em.To = []string{"example@163.com"}
	em.Subject = "测试"
	em.HTML = []byte("测试")

	// 发送邮件,超时20秒
	err = pool.Send(em, time.Second*20)

	// 断言发送是否错误
	assert.NoError(t, err)

	// 等待5秒
	time.Sleep(time.Second * 5)
}
