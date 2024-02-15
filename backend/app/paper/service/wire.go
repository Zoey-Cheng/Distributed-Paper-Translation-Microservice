//go:build wireinject
// +build wireinject

package main

import (
	v1 "paper-translation/api/paper/service/v1"
	"paper-translation/app/paper/service/paper"
	"paper-translation/pkg/ds"
	"paper-translation/pkg/service"

	"github.com/google/wire"
	"go-micro.dev/v4"
)

func InitApp() micro.Service {
	panic(wire.Build(
		service.ProviderSet,
		ds.NewMongoClient,
		ds.NewMongoDatabase,
		paper.NewMongoPaperRepository, wire.Bind(new(paper.PaperRepository), new(*paper.MongoPaperRepository)),
		NewFileService,
		NewOCRService,
		NewTranslationService,
		NewEmailService,
		paper.NewPaperService, wire.Bind(new(v1.PaperServiceHandler), new(*paper.PaperService)),
		NewService,
	))
}
