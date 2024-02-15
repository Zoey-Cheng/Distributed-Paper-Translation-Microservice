package main

import (
	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
	v1 "paper-translation/api/translation/service/v1"
	"paper-translation/pkg/errutil"
	"paper-translation/pkg/service"
)

func NewService(registry registry.Registry, config config.Config, handler v1.TranslationServiceHandler) micro.Service {
	svc := micro.NewService(
		micro.Name(service.TranslationServiceName),
		micro.Address(config.Get("server", "addr").String(":4000")),
		micro.Registry(registry),
	)
	err := v1.RegisterTranslationServiceHandler(svc.Server(), handler)
	errutil.PanicIfErr(err)
	svc.Init()
	return svc
}
