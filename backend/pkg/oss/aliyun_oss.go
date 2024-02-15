package oss

import (
	"fmt"
	"paper-translation/pkg/errutil" // 错误处理

	"github.com/aliyun/aliyun-oss-go-sdk/oss" // 阿里云OSS SDK
)

/**
* 获取阿里云OSS客户端
* @param region - 区域
* @param keyID - AccessKey ID
* @param secret - AccessKey Secret
* @return OSS客户端
 */
func MustGetAliYunOSS(region string, keyID string, secret string) *oss.Client {

	endpoint := fmt.Sprintf("https://oss-%s.aliyuncs.com", region) // 构建endpoint

	client, err := oss.New(endpoint, keyID, secret) // 创建客户端

	errutil.PanicIfErr(err) // 错误处理

	return client
}
