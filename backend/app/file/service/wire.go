//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-micro.dev/v4"
	"paper-translation/api/file/service/v1"
	"paper-translation/app/file/service/file"
	"paper-translation/pkg/ds"
	"paper-translation/pkg/service"
)

func InitApp() micro.Service {
	panic(wire.Build(
		service.ProviderSet,
		ds.NewMongoClient,
		ds.NewRedisClient,
		ds.NewMongoDatabase,
		file.NewMongoFileRepository,
		file.NewFileService, wire.Bind(new(v1.FileServiceHandler), new(*file.FileService)),
		NewService,
	))
}
