package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss" // OSS SDK
	"go-micro.dev/v4/config"                  // 配置管理器
)

/**
* 从配置中获取阿里云OSS客户端
* @param config - 配置管理器
* @return OSS客户端
 */
func NewAliYunOSS(config config.Config) *oss.Client {

	region := config.Get("aliyun", "oss", "region").String("cn-beijing") // 获取region
	keyID := config.Get("aliyun", "oss", "key_id").String("default")     // 获取key id
	secret := config.Get("aliyun", "oss", "secret").String("default")    // 获取secret

	return MustGetAliYunOSS(region, keyID, secret) // 获取客户端
}
