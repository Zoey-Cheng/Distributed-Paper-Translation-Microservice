//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-micro.dev/v4"
	v1 "paper-translation/api/translation/service/v1"
	"paper-translation/app/translation/service/translation"
	"paper-translation/pkg/ds"
	"paper-translation/pkg/service"
	"paper-translation/pkg/signal"
	xfspark "paper-translation/pkg/xf-spark"
)

func InitApp() micro.Service {
	panic(wire.Build(
		service.ProviderSet,
		xfspark.NewXFSpark,
		signal.NewSignalFactory,
		ds.NewRedisClient,
		translation.NewTranslationService, wire.Bind(new(v1.TranslationServiceHandler), new(*translation.TranslationService)),
		NewService,
	))
}
