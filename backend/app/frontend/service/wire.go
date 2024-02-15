//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-micro.dev/v4/web"
	"paper-translation/pkg/oss"
	"paper-translation/pkg/service"
)

func InitApp() web.Service {
	panic(wire.Build(
		service.ProviderSet,
		oss.NewAliYunOSS,
		NewFileService,
		NewPaperService,
		NewRoute,
		NewService,
	))
}
