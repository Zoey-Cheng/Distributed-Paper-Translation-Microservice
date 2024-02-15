package email

import (
	"bytes"
	"context"
	"html/template"
	"log"
	v1 "paper-translation/api/email/service/v1"
	"time"

	"github.com/jordan-wright/email"
	"go-micro.dev/v4/config"
)

// EmailService 是电子邮件服务的实现。
type EmailService struct {
	queue     chan *v1.SendEmailParam // 电子邮件发送队列
	emailPool *email.Pool             // 电子邮件池
	config    config.Config           // 配置信息
}

// NewEmailService 创建一个新的 EmailService 实例。
//
// 参数:
// - emailPool (*email.Pool): 电子邮件池。
// - config (config.Config): 配置信息。
//
// 返回值:
// - *EmailService: EmailService 实例。
func NewEmailService(emailPool *email.Pool, config config.Config) *EmailService {
	serv := &EmailService{emailPool: emailPool, config: config}
	serv.Start()
	return serv
}

// Start 启动电子邮件服务，监听发送队列并处理电子邮件发送请求。
func (serv *EmailService) Start() {
	serv.queue = make(chan *v1.SendEmailParam, 1000)
	go func() {
		for param := range serv.queue {
			tmpl, err := template.New("email").Parse(param.Template)
			if err != nil {
				log.Printf("parse template err: %+v", err)
				continue
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, &param.Vars)
			if err != nil {
				log.Printf("execute template err: %+v", err)
				continue
			}

			em := email.NewEmail()
			em.From = config.Get("email", "username").String("golangproject@163.com")
			em.To = []string{param.EmailTo}
			em.Subject = param.Subject
			em.HTML = buf.Bytes()

			err = serv.emailPool.Send(em, 10*time.Second)
			if err != nil {
				log.Printf("send email err: %+v", err)
				continue
			}
		}
	}()
}

// SendEmail 发送电子邮件。
//
// 参数:
// - ctx (context.Context): 上下文。
// - param (*v1.SendEmailParam): 发送电子邮件的参数。
// - status (*v1.EmailStatus): 电子邮件发送状态。
//
// 返回值:
// - error: 错误信息，如果发生错误。
func (serv *EmailService) SendEmail(ctx context.Context, param *v1.SendEmailParam, status *v1.EmailStatus) error {
	serv.queue <- param
	status.Status = true
	return nil
}
