//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-micro.dev/v4"
	v1 "paper-translation/api/email/service/v1"
	"paper-translation/app/email/service/email"
	emailPool "paper-translation/pkg/email"
	"paper-translation/pkg/service"
)

func InitApp() micro.Service {
	panic(wire.Build(
		service.ProviderSet,
		emailPool.NewEmailPool,
		email.NewEmailService, wire.Bind(new(v1.EmailServiceHandler), new(*email.EmailService)),
		NewService,
	))
}
