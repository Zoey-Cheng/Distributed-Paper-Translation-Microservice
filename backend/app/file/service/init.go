package main

import (
	v1 "paper-translation/api/file/service/v1"
	"paper-translation/pkg/errutil"
	"paper-translation/pkg/service"

	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
)

// NewService 创建并返回一个新的 Micro 服务实例。
//
// 参数:
// - registry registry.Registry: Micro 服务注册表。
// - config config.Config: 配置信息。
// - handler v1.FileServiceHandler: 文件服务处理器。
//
// 返回值:
// - micro.Service: 创建的 Micro 服务实例。
func NewService(registry registry.Registry, config config.Config, handler v1.FileServiceHandler) micro.Service {
	svc := micro.NewService(
		micro.Name(service.FileServiceName),                         // 设置服务名称
		micro.Address(config.Get("server", "addr").String(":4000")), // 设置服务地址
		micro.Registry(registry),                                    // 设置服务注册表
	)
	err := v1.RegisterFileServiceHandler(svc.Server(), handler) // 注册文件服务处理器
	errutil.PanicIfErr(err)                                     // 如果有错误，立即中断程序
	svc.Init()                                                  // 初始化服务
	return svc                                                  // 返回创建的 Micro 服务实例
}
