package main

import (
	v1 "paper-translation/api/ocr/service/v1"
	"paper-translation/pkg/errutil"
	"paper-translation/pkg/service"

	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
)

// NewService 创建并返回一个新的 OCR 微服务实例。
//
// 参数:
// - registry registry.Registry: 微服务注册表。
// - config config.Config: 配置信息。
// - handler v1.OCRServiceHandler: OCR 微服务的处理程序。
//
// 返回值:
// - micro.Service: 创建的 OCR 微服务实例。
func NewService(registry registry.Registry, config config.Config, handler v1.OCRServiceHandler) micro.Service {
	svc := micro.NewService(
		micro.Name(service.OCRServiceName),                          // 设置服务名称
		micro.Address(config.Get("server", "addr").String(":4000")), // 设置服务地址
		micro.Registry(registry),                                    // 设置微服务注册表
	)
	err := v1.RegisterOCRServiceHandler(svc.Server(), handler) // 注册 OCR 微服务处理程序
	errutil.PanicIfErr(err)                                    // 如果发生错误，触发恐慌
	svc.Init()                                                 // 初始化服务
	return svc                                                 // 返回创建的 OCR 微服务实例
}
