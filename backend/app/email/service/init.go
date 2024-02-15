package main

import (
	v1 "paper-translation/api/email/service/v1"
	"paper-translation/pkg/errutil"
	"paper-translation/pkg/service"

	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
)

// NewService 创建一个新的微服务实例。
//
// 参数:
// - registry (registry.Registry): 注册中心。
// - config (config.Config): 配置信息。
// - handler (v1.EmailServiceHandler): EmailService 的处理程序。
//
// 返回值:
// - micro.Service: 微服务实例。
func NewService(registry registry.Registry, config config.Config, handler v1.EmailServiceHandler) micro.Service {
	// 创建微服务实例并配置基本属性。
	svc := micro.NewService(
		micro.Name(service.EmailServiceName),
		micro.Address(config.Get("server", "addr").String(":4000")),
		micro.Registry(registry),
	)

	// 注册 EmailService 的处理程序。
	err := v1.RegisterEmailServiceHandler(svc.Server(), handler)
	errutil.PanicIfErr(err)

	// 初始化微服务。
	svc.Init()
	return svc
}
