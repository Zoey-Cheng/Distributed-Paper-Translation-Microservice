package service

import (
	"os"
)

var (
	// EtcdAddr 存储用于连接到 etcd 服务的地址。
	EtcdAddr string
	// ConfigPrefix 存储配置的前缀，用于在 etcd 中查找配置。
	ConfigPrefix string
)

func init() {
	// 从环境变量 ETCD_ADDR 获取 etcd 服务地址，如果未设置，则使用默认值 "localhost:2379"。
	EtcdAddr = os.Getenv("ETCD_ADDR")
	if EtcdAddr == "" {
		EtcdAddr = "localhost:2379"
	}

	// 从环境变量 CONFIG_PREFIX 获取配置前缀，如果未设置，则使用默认值 "/config/default"。
	ConfigPrefix = os.Getenv("CONFIG_PREFIX")
	if ConfigPrefix == "" {
		ConfigPrefix = "/config/default"
	}
}
