//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-micro.dev/v4"
	v1 "paper-translation/api/ocr/service/v1"
	"paper-translation/app/ocr/service/ocr"
	"paper-translation/pkg/ds"
	aliYunOCR "paper-translation/pkg/ocr"
	"paper-translation/pkg/oss"
	"paper-translation/pkg/service"
)

func InitApp() micro.Service {
	panic(wire.Build(
		service.ProviderSet,
		ds.NewMongoClient,
		ds.NewMongoDatabase,
		ds.NewRedisClient,
		aliYunOCR.NewAliYunOCR,
		oss.NewAliYunOSS,
		ocr.NewMongoOCRRepository, wire.Bind(new(ocr.OCRRepository), new(*ocr.MongoOCRRepository)),
		ocr.NewOCRService, wire.Bind(new(v1.OCRServiceHandler), new(*ocr.OCRService)),
		NewService,
	))
}
