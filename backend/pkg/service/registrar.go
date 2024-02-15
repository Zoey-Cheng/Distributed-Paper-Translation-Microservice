package service

import (
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
)

/**
 * NewRegistry 函数创建并返回一个新的服务注册表对象，用于服务发现和注册。
 *
 * @return registry.Registry - 服务注册表对象
 */
func NewRegistry() registry.Registry {
	// 创建一个 etcd 服务注册表，使用 EtcdAddr 作为地址，不启用安全连接。
	return etcd.NewRegistry(
		registry.Addrs(EtcdAddr),
		registry.Secure(false),
	)
}
