package service

import (
	"log"
	"paper-translation/pkg/errutil"

	"github.com/go-micro/plugins/v4/config/source/etcd"
	"go-micro.dev/v4/config"
)

/**
 * NewConfig 函数创建并返回一个用于配置管理的新配置对象。
 *
 * @return config.Config - 配置对象
 */
func NewConfig() config.Config {
	// 创建一个新的配置对象
	cfg, err := config.NewConfig()
	errutil.PanicIfErr(err)

	// 打印日志，指示从 etcd 读取配置
	log.Printf("start read config from etcd: %s, prefix: %s", EtcdAddr, ConfigPrefix)

	// 从 etcd 加载配置
	err = cfg.Load(etcd.NewSource(
		etcd.WithAddress(EtcdAddr),
		etcd.WithPrefix(ConfigPrefix),
		etcd.StripPrefix(true),
	))
	errutil.PanicIfErr(err)

	// 打印加载的配置
	log.Printf("load config is: %+v", cfg.Map())

	// 返回配置对象
	return cfg
}
