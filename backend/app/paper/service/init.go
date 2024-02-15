package main

import (
	"go-micro.dev/v4/client"
	es "paper-translation/api/email/service/v1"
	fs "paper-translation/api/file/service/v1"
	os "paper-translation/api/ocr/service/v1"
	v1 "paper-translation/api/paper/service/v1"
	ts "paper-translation/api/translation/service/v1"
	"paper-translation/pkg/errutil"
	"paper-translation/pkg/service"
	"time"

	"go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
)

func NewService(registry registry.Registry, config config.Config, handler v1.PaperServiceHandler) micro.Service {
	svc := micro.NewService(
		micro.Name(service.PaperServiceName),
		micro.Address(config.Get("server", "addr").String(":4000")),
		micro.Registry(registry),
	)
	err := v1.RegisterPaperServiceHandler(svc.Server(), handler)
	errutil.PanicIfErr(err)
	svc.Init()
	return svc
}

func NewFileService(registry registry.Registry) fs.FileService {
	cli := client.NewClient(
		client.Registry(registry),
		client.RequestTimeout(time.Second*300),
		client.RequestTimeout(time.Second*300),
		client.StreamTimeout(time.Second*300),
	)
	return fs.NewFileService(service.FileServiceName, cli)
}

func NewOCRService(registry registry.Registry) os.OCRService {
	cli := client.NewClient(
		client.Registry(registry),
	)
	return os.NewOCRService(service.OCRServiceName, cli)
}

func NewTranslationService(registry registry.Registry) ts.TranslationService {
	cli := client.NewClient(
		client.Registry(registry),
	)
	return ts.NewTranslationService(service.TranslationServiceName, cli)
}

func NewEmailService(registry registry.Registry) es.EmailService {
	cli := client.NewClient(
		client.Registry(registry),
	)
	return es.NewEmailService(service.EmailServiceName, cli)
}
